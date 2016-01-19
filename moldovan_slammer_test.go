package main

import "testing"

func TestBuildSQL(t *testing.T) {
	template := "INSERT INTO `floop` VALUES ('{guid}','{guid:0}',{int:-2000:0},{int:100:1000},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL)"
	_, err := buildSQL(template)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkBuildSQL(b *testing.B) {
	template := "INSERT INTO `floop` VALUES ('{guid}','{guid:0}',{int:-2000:0},{int:100:1000},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL)"

	for n := 0; n < b.N; n++ {
		buildSQL(template)
	}
}
