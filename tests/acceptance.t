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

# make sure we've got `curl`
{
  my $shell_ok = system('which curl > /dev/null');
  is( $shell_ok, 0, 'can find `curl` in current path' );
}

# the service counts words, normalizing to lower-case
{
  my $filename = 'quirky-contemplating';
  my $dict = TestHelpers->POST(port => $port, filename => $filename, content => <<END);
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

  my $cached = TestHelpers->POST(port => $port, filename => $filename, content => 'gets ignored');
  is_deeply( $cached, $dict, 'content is cached based on filename' );
}

# we can GET and DELETE from the cache
{
  # sanity check: this file isn't in the service yet
  my $test_file = 'bicyclists-bellicose';
  my $cached_files = TestHelpers->GET(port => $port);
  is_deeply(
    [grep { $_ eq $test_file } @$cached_files],
    [],
    "$test_file not in cache"
  );

  # stick it in the service and try again
  TestHelpers->POST(port => $port, filename => $test_file, content => 'unmindful-dervishes');
  $cached_files = TestHelpers->GET(port => $port);
  is_deeply(
    [grep { $_ eq $test_file } @$cached_files],
    [$test_file],
    "$test_file in cache"
  );

  # delete it out, yo
  TestHelpers->DELETE(port => $port, filename => $test_file);
  $cached_files = TestHelpers->GET(port => $port);
  is_deeply(
    [grep { $_ eq $test_file } @$cached_files],
    [],
    "$test_file no longer in cache"
  );
}
