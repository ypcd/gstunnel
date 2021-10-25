module gstunnellib

go 1.17

require (
	google.golang.org/protobuf v1.27.1
	//gstunnellib/gspackoper v0.0.0-00010101000000-000000000000
	//gstunnellib/gsrand v0.0.0-00010101000000-000000000000
	timerm v0.0.0-00010101000000-000000000000
)

replace timerm => ../timerm

//replace gstunnellib/gspackoper => ./gspackoper

//replace gstunnellib/gsrand => ./gsrand
