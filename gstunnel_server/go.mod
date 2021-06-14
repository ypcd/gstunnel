module gstunnel_server_mod

go 1.16

require (
	gstunnellib v0.0.0-00010101000000-000000000000
	timerm v0.0.0-00010101000000-000000000000
)

replace gstunnellib => ../gstunnellib

replace timerm => ../timerm
