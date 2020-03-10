# Gravity Data Source Adapter

Data source adapter offers a way for recevieing events and data from external system.

## Update proto definition

Rebuild to apply `proto` changes, just run commands below:

```shell
cd pb
protoc --go_out=plugins=grpc:. *.proto
```

## License

Licensed under the MIT License

## Authors

Copyright(c) 2020 Fred Chien <<fred@brobridge.com>>
