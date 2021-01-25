run:
	#Init main service
	cd api-mgr/cmd; go run main.go &

	# Run service A
	cd lib/examples/service-A; go run sample.go &

	# Run service B
	cd lib/examples/service-B; go run sample.go &

	# Run service C
	cd lib/examples/service-C; go run sample.go &

	# Run Web service
	cd web; npm install; npm run dev &