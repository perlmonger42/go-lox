#!/usr/bin/env perl
use strict;
use warnings;

my $line;
my $label;

while ($line = <>) {
  $label = $1 if $line =~ /^--- FAIL: (\w+) \(.*\)/;
  last if $line =~ /^got:\s*$/;
}

open (my $got, '>', ".test-got.txt") or
  die "Couldn't open .test-got.txt";
print $got "===== $label =====\n";
while ($line = <>) {
  last if $line =~ /^want:\s*$/;
  print $got $line;
}
close $got;


open (my $want, '>', ".test-want.txt") or
  die "Couldn't open .test-want.txt";
print $want "===== $label =====\n";
while ($line = <>) {
  last if $line =~ /^FAIL\s*$/;
  print $want $line;
}
close $want;

system "opendiff", ".test-got.txt", ".test-want.txt";
# vim: filetype=perl shiftwidth=8 tabstop=8 noexpandtab
