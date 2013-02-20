bitly API golang library
========================

## Run tests

Your username is the lowercase name shown when you login to bitly, your access token can be fetched using the following (http://dev.bitly.com/authentication.html ):

    curl -u "username:password" -X POST "https://api-ssl.bitly.com/oauth/access_token"

To run the tests either export the environment variable or set it up inline before calling `nosetests`:

    BITLY_API_TOKEN=<accesstoken> go test

## API Documentation

http://dev.bitly.com/
