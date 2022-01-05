module main

go 1.17

replace (
	lin_common => ../lin_common
	lin_cor_pool => ../lin_cor_pool
)

require lin_cor_pool v0.0.0-00010101000000-000000000000

require lin_common v0.0.0-00010101000000-000000000000 // indirect
