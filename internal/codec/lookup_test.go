package codec

import "testing"

func TestByName(t *testing.T) {
	c, ok := ByName("zstd")
	if !ok || c.Bin != "zstd" {
		t.Fatal("zstd not found")
	}
	if _, ok := ByName("nope"); ok {
		t.Fatal("unknown matched")
	}
}

func TestByExtensionLongestMatch(t *testing.T) {
	c, ok := ByExtension("archive.tar.gz")
	if !ok || c.Name != "gzip" {
		t.Fatalf("want gzip, got %v ok=%v", c.Name, ok)
	}
	c, ok = ByExtension("plain.tar")
	if !ok || c.Name != "tar" {
		t.Fatalf("want tar, got %v", c.Name)
	}
	if _, ok := ByExtension("file.txt"); ok {
		t.Fatal(".txt must not match")
	}
}
