use strict;
use warnings;
package TestHelpers;
use base qw(Exporter);

use JSON;

# "You could use a CPAN library for a lot of this."
# Yep, that's true - but for a one-off project, I don't
# see much value in forcing the end user to install a bunch
# of CPAN libraries. So I'm coding against a possibly old core Perl.

my $PORT;
sub set_port { $PORT = pop }

sub wait_for_server_to_start_running {
  my ($class, $port) = @_;

  my $too_many;

  # Just use 'lsof' to check the port every half second.
  while (0 != system("lsof -i :$port > /dev/null")) {
    select undef, undef, undef, 0.5;
    if ($too_many++ > 20) {
      die("Gave up after " . ($too_many / 2) . " seconds waiting for port :$port to be attached");
    }
  }

  return;
}

sub _simple_curl {
  my ($class, %args) = @_;

  my $endpoint = $args{endpoint} || '/files';

  my $request = "http://localhost:${PORT}${endpoint}";
  $request .= "?filename=$args{filename}" if exists $args{filename};

  chomp(my $response = `curl -s -X $args{method} '$request'`);
  my $json = JSON->new->allow_nonref;
  return $json->decode($response);
}

sub DELETE {
  my ($class, %args) = @_;
  return $class->_simple_curl(method => 'DELETE', %args);
}

sub GET {
  my ($class, %args) = @_;
  return $class->_simple_curl(method => 'GET', %args);
}

sub PUT {
  my ($class, %args) = @_;
  return $class->_simple_curl(method => 'PUT', %args);
}

sub POST {
  my ($filename, %args) = @_;

  my $endpoint = $args{endpoint} || '/files';

  # bi-directional pipes are a nuisance in perl, just pipe $content to curl and send output to a file
  my $tempfile = "/tmp/output.$$";
  open my $fh, '+>', $tempfile or die "write $tempfile: $!\n";
  print $fh $args{content};
  close $fh;
  my $json = `curl -s -F filename=$args{filename} -F file=\@${tempfile} 'http://localhost:${PORT}${endpoint}'`;

  unlink $tempfile;

  # convert json to a data structure
  return decode_json($json);
}

1;
