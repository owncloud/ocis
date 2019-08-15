# reva-hyper

## Organizational structure
* using reva as library
* building services around the library functionalities
* standalone repos for each service
* hyperreva as wrapper repository for standalone binary

Suggested repos naming schema:
* reva-hyper 
  * reva-phoenix
  * reva-my-service1
  * reva-my-service2
  * ...

## Technical details
* microservices: [go-kit/kit](https://github.com/go-kit/kit)
* cli: [spf13/cobra](https://github.com/spf13/cobra)
* configuration: [spf13/viper](https://github.com/spf13/viper)

## Next goals
* [ ] working integration demo
* [ ] minimal drone pipeline for nightly builds
