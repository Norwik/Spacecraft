<p align="center">
    <img src="/assets/logo.png?v=1.0.1" width="200" />
    <h3 align="center">Spacecraft</h3>
    <p align="center">Helmet getting started example with docker compose.</p>
    <p align="center">
        <a href="https://github.com/spacewalkio/spacecraft/blob/master/LICENSE">
            <img src="https://img.shields.io/badge/LICENSE-MIT-E74C3C.svg">
        </a>
    </p>
</p>


## Documentation

Login into github packages registry to be able to pull helmet docker image.

```zsh
$ apt install pass jq

# Get the token from https://github.com/settings/tokens
$ export GITHUB_TOKEN="~token~goes~here~"
$ echo $GITHUB_TOKEN | docker login https://docker.pkg.github.com -u username --password-stdin
```

Run the example with the following command:

```zsh
$ docker-compose up -d
```

Create an auth method for internal communication (orders microservice ----> customers microservice)

```zsh
$ curl -X POST \
    http://127.0.0.1:8000/apigw/api/v1/auth/method \
    -H 'X-API-KEY: 6c68b836-6f8e-465e-b59f-89c1db53afca' \
    -d '{"name": "internal_oauth_method","type": "oauth_authentication","description": "Internal Microservices Communication","endpoints": "orders_service;customers_service"}' | jq .

{
  "id": 1,
  "name": "internal_oauth_method",
  "description": "Internal Microservices Communication",
  "type": "oauth_authentication",
  "endpoints": "orders_service;customers_service",
  "createdAt": "2021-09-28T17:30:19.96193235Z",
  "updatedAt": "2021-09-28T17:30:19.96193235Z"
}
```

Create an auth method for external communication (end user --> orders microservice)

```zsh
$ curl -X POST \
    http://127.0.0.1:8000/apigw/api/v1/auth/method \
    -H 'X-API-KEY: 6c68b836-6f8e-465e-b59f-89c1db53afca' \
    -d '{"name": "external_key_method","type": "key_authentication","description": "External Microservices Communication","endpoints": "orders_service"}' | jq .

{
  "id": 2,
  "name": "external_key_method",
  "description": "External Microservices Communication",
  "type": "key_authentication",
  "endpoints": "orders_service",
  "createdAt": "2021-09-28T17:31:26.541449304Z",
  "updatedAt": "2021-09-28T17:31:26.541449304Z"
}
```

Create Oauth credentials for orders microservice

```zsh
$ curl -X POST \
    http://127.0.0.1:8000/apigw/api/v1/auth/oauth \
    -H 'X-API-KEY: 6c68b836-6f8e-465e-b59f-89c1db53afca' \
    -d '{"name": "orders_microservice","clientID": "orders_microservice.spacewalkio","clientSecret": "4a0d4910-e902-432e-9f53-fad719a3d224","meta": "name=orders;entity=emea","authMethodID":1}' | jq .

{
  "id": 1,
  "name": "orders_microservice",
  "clientID": "orders_microservice.spacewalkio",
  "clientSecret": "4a0d4910-e902-432e-9f53-fad719a3d224",
  "meta": "name=orders;entity=emea",
  "authMethodID": 1,
  "createdAt": "2021-09-28T17:41:47.387256326Z",
  "updatedAt": "2021-09-28T17:41:47.387256326Z"
}
```

Create an API Key for end user to call the orders microservice

```zsh
$ curl -X POST \
    http://127.0.0.1:8000/apigw/api/v1/auth/key \
    -H 'X-API-KEY: 6c68b836-6f8e-465e-b59f-89c1db53afca' \
    -d '{"name": "j.doe","apiKey": "2f521881-7481-412f-9948-b466498b59d3","meta": "id=2001;entity=emea","authMethodID": 2}' | jq .

{
  "id": 1,
  "name": "j.doe",
  "apiKey": "2f521881-7481-412f-9948-b466498b59d3",
  "meta": "id=2001;entity=emea",
  "authMethodID": 2,
  "createdAt": "2021-09-28T17:42:01.035173754Z",
  "updatedAt": "2021-09-28T17:42:01.035173754Z"
}
```

Start calling the orders microservice using the recently created API Key.

```zsh
$ curl -X GET http://127.0.0.1:8000/orders/v1/order/1 -H 'X-API-KEY: 2f521881-7481-412f-9948-b466498b59d3' -v
```

Now you can visit [Prometheus](http://127.0.0.1:9090/targets) and [Grafana dashboards](http://127.0.0.1:3000/) (Login admin/admin). Add a new data source to grafana (http://prometheus:9090) Then load the `grafana_dashboard.json` file to get a dashboard like `grafana_screenshot01.png`.

Stop the running containers with the following command:

```zsh
$ docker-compose down
```


## Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, Spacecraft is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/spacewalkio/spacecraft/releases) for changelogs for each release version of Spacecraft. It contains summaries of the most noteworthy changes made in each release.


## Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/spacewalkio/spacecraft/issues


## Security Issues

If you discover a security vulnerability within Spacecraft, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


## Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


## License

© 2021, SpaceWalk. Released under [MIT License](https://opensource.org/licenses/mit-license.php).

**Spacecraft** is authored and maintained by [@SpaceWalk](http://github.com/spacewalkio).
