package moldova

import "testing"

// TODO Test each random function individually, under a number of inputs to make supported
// all the options behave as expected.

func TestBuildSQL(t *testing.T) {
	template := "INSERT INTO floof VALUES ('{guid}','{guid:0}','{country}',{int:-2000:0},{int:100:1000},{float:-1000.0:-540.0},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL,-3)"
	_, err := ParseTemplate(template)
	if err != nil {
		t.Error(err)
	}
}

func TestCountries(t *testing.T) {
	template := "INSERT INTO `floop` VALUES ('{country}','{country:up:0}','{country}','{country:down:1}')"
	_, err := ParseTemplate(template)
	if err != nil {
		t.Error(err)
	}
}

func BenchmarkBuildSQL(b *testing.B) {
	template := "INSERT INTO `floop` VALUES ('{guid}','{guid:0}',{int:-2000:0},{int:100:1000},{int:1:40},'{now}','{now:0}','{char:2:up}',NULL)"

	for n := 0; n < b.N; n++ {
		ParseTemplate(template)
	}
}
