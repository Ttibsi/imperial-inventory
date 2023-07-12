# Imperial Inventory

You are R3-D3 and were just appointed the general of the imperial fleet. Your
first action as the new general is to digitalise the imperial fleet inventory.

### TO USE
Prerequisites: Docker

In one terminal window, run `make db` to start the database container.
In another terminal window, run `go run .` to start up the webserver.

### For Future Improvements
* There are no unit testing here as I've run out of time -- I'd like to include
unit tests on each endpoint. This should be simple to do, using `httptest.NewRecorder(()`
to track the output when a request is made and checking the body response there.

* Data sanitisation -- there's a lot of data coming in from the user on the PUT
and POST requests, and no validation at all. The little sanitisation there is
currently is given by default by `database/sql`.

* Docker volumes for database persistance -- The biggest challenge here was
setting up the database as I've not used mySQL specifically before. Currently,
the database is recreated from scratch each time the container is started. Given
enough time, I would have instead set up a docker volume so the data is stored
on the host machine instead.
