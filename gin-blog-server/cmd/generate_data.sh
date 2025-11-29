cd generate-data

# all | config | auth | page | resource
# all Generate all information
# config Generate configuration information
# auth Generate default roles admin, guest, and corresponding default users admin, guest
# page Generate default page information
# resource Generate default resource information
go run main.go -t "all"