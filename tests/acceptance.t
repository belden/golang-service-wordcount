#!/usr/bin/env perl

use strict;
use warnings;
use Test::More; END { done_testing }

use Getopt::Long;
use FindBin;
use lib "$FindBin::Bin/";
use TestHelpers;

GetOptions(
  'port=i' => \(my $port = 3000),
);

TestHelpers->wait_for_server_to_start_running($port);
TestHelpers->set_port($port);

# make sure we've got `curl`
{
  my $shell_ok = system('which curl > /dev/null');
  is( $shell_ok, 0, 'can find `curl` in current path' );
}

# the service counts words, normalizing to lower-case
{
  my $filename = 'quirky-contemplating';
  my $dict = TestHelpers->POST(filename => $filename, content => <<END);
This is <a> sentence with <a> few extra characters.
You can't forget that an apostrophe is a legitimate character for English words.
And though underscores aren't legitimate in English, some mothers-in-law love hyphens.
The reason for this is unknown. Is it? IT IS.
END
  my $expected = +{
    Total => 43,
    Words => {
      is               => 5,
      a                => 3,
      english          => 2,
      for              => 2,
      it               => 2,
      legitimate       => 2,
      this             => 2,
      an               => 1,
      and              => 1,
      apostrophe       => 1,
      "aren't"         => 1,
      "can't"          => 1,
      character        => 1,
      characters       => 1,
      extra            => 1,
      few              => 1,
      forget           => 1,
      hyphens          => 1,
      in               => 1,
      love             => 1,
      'mothers-in-law' => 1,
      reason           => 1,
      sentence         => 1,
      some             => 1,
      that             => 1,
      the              => 1,
      though           => 1,
      underscores      => 1,
      unknown          => 1,
      with             => 1,
      words            => 1,
      you              => 1,
    },
  };

  is_deeply( $dict, $expected, 'got right data for test text' );

  my $cached = TestHelpers->POST(filename => $filename, content => 'gets ignored');
  is_deeply( $cached, $dict, 'content is cached based on filename' );
}

# we can GET and DELETE from the cache
{
  # sanity check: this file isn't in the service yet
  my $test_file = 'bicyclists-bellicose';
  my $cached_files = TestHelpers->GET;
  is_deeply(
    [grep { $_ eq $test_file } @$cached_files],
    [],
    "$test_file not in cache"
  );

  # stick it in the service and try again
  TestHelpers->POST(filename => $test_file, content => 'unmindful-dervishes');
  $cached_files = TestHelpers->GET;
  is_deeply(
    [grep { $_ eq $test_file } @$cached_files],
    [$test_file],
    "$test_file in cache"
  );

  # delete it out, yo
  TestHelpers->DELETE(filename => $test_file);
  $cached_files = TestHelpers->GET;
  is_deeply(
    [grep { $_ eq $test_file } @$cached_files],
    [],
    "$test_file no longer in cache"
  );
}

# we can GET a specific file
{
	my $filename = 'Garry-colors';
	my $post = TestHelpers->POST(filename => $filename, content => 'deadpanned Porter songwriter-denote');
	my $get = TestHelpers->GET(filename => $filename);
	is_deeply( $get, $post, 'we can retrieve from cache' );
}

# /admin/ routes operate on all files
{
  # some sanity setup: stick something in the cache and ensure it's listed as a file we can GET
  my $filename = 'casuist-monied';
  TestHelpers->POST(filename => $filename, content => "this is ${filename} content");
  my $cached_files = TestHelpers->GET();
  ok( 1 == scalar(grep { $_ eq $filename } @$cached_files), "${filename} shows up in cache" );

  # let's bulk-DELETE everything and ensure our canary file is gone
	TestHelpers->DELETE(endpoint => '/admin/files');
  is_deeply( TestHelpers->GET, [], 'DELETE /admin/files deletes files' );

  # Insert some files to the cache
  my @stub_files = qw(molecules-Clayton bifurcated-unmanned);
	my @expected_json = map {
		TestHelpers->POST(filename => $_, content => "this is $_");
	} @stub_files;

  # bulk-fetch all files
  my $bulk_get = TestHelpers->GET(endpoint => '/admin/files');
  is_deeply(
    $bulk_get,
    {
      map { $_ => +{Total => 3, Words => +{this => 1, is => 1, lc($_) => 1}} }
      @stub_files
    },
    'GET /admin/files returns expected data'
  );
}

sub show {
  use Data::Dumper;
  $Data::Dumper::SortKeys = $Data::Dumper::SortKeys = 1;
  print STDERR Dumper(shift);
}

# multi-GET will sum things up for us
{
  my @files = qw(clangs-sweetening advising-Heliopolis succotashs-togas);
  foreach my $file (@files) {
    TestHelpers->POST(filename => $file, content => "to be or not to be, ${file} is the question");
  }
  my $multi = TestHelpers->GET(endpoint => "/files?filename=$files[0]&filename=$files[2]");
  is_deeply( $multi, +{
    Total => 20, # because the test text for each "file" is 10 words long
    Words => +{
      to => 4,
      be => 4,
      or => 2,
      not => 2,
      lc($files[0]) => 1,
      lc($files[2]) => 1,
      is => 2,
      the => 2,
      question => 2,
    },
  }, 'GET /files?filename=foo&filename=bar sums things up' ) or show($multi);
}
