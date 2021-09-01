module example.com/application

go 1.13

replace example.com/search => ./search

require (
	example.com/matchers v0.0.0-00010101000000-000000000000
	example.com/search v0.0.0-00010101000000-000000000000
)

replace example.com/matchers => ./matchers
