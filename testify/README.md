Testify - Thou Shalt Write Tests
================================

ℹ️ We are working on testify v2 and would love to hear what you'd like to see in it, have your say
here: https://cutt.ly/testify

[![Build Status](https://travis-ci.org/stretchr/testify.svg)](https://travis-ci.org/stretchr/testify) [![Go Report Card](https://goreportcard.com/badge/github.com/stretchr/testify)](https://goreportcard.com/report/github.com/stretchr/testify) [![PkgGoDev](https://pkg.go.dev/badge/github.com/stretchr/testify)](https://pkg.go.dev/github.com/stretchr/testify)

Go code (golang) set of packages that provide many tools for testifying that your code will behave as you intend.

Features include:

* [Easy assertions](#assert-package)
* [Mocking](#mock-package)
* [Testing suite interfaces and functions](#suite-package)

Get started:

* Install testify with [one line of code](#installation), or [update it with another](#staying-up-to-date)
* For an introduction to writing test code in Go, see http://golang.org/doc/code.html#Testing
* Check out the API Documentation http://godoc.org/github.com/stretchr/testify
* To make your testing life easier, check out our other project, [gorc](http://github.com/stretchr/gorc)
* A little about [Test-Driven Development (TDD)](http://en.wikipedia.org/wiki/Test-driven_development)

[`assert`](http://godoc.org/github.com/stretchr/testify/assert "API documentation") package
-------------------------------------------------------------------------------------------

The `assert` package provides some helpful methods that allow you to write better test code in Go.

* Prints friendly, easy to read failure descriptions
* Allows for very readable code
* Optionally annotate each assertion with a message

See it in action:

```go
package yours

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

	// assert equality
	assert.Equal(t, 123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(t, 123, 456, "they should not be equal")

	// assert for nil (good for errors)
	assert.Nil(t, object)

	// assert for not nil (good when you expect something)
	if assert.NotNil(t, object) {

		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		assert.Equal(t, "Something", object.Value)

	}

}
```

* Every assert func takes the `testing.T` object as the first argument. This is how it writes the errors out through the
  normal `go test` capabilities.
* Every assert func returns a bool indicating whether the assertion was successful or not, this is useful for if you
  want to go on making further assertions under certain conditions.

if you assert many times, use the below:

```go
package yours

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {
	assert := assert.New(t)

	// assert equality
	assert.Equal(123, 123, "they should be equal")

	// assert inequality
	assert.NotEqual(123, 456, "they should not be equal")

	// assert for nil (good for errors)
	assert.Nil(object)

	// assert for not nil (good when you expect something)
	if assert.NotNil(object) {

		// now we know that object isn't nil, we are safe to make
		// further assertions without causing any errors
		assert.Equal("Something", object.Value)
	}
}
```

[`require`](http://godoc.org/github.com/stretchr/testify/require "API documentation") package
---------------------------------------------------------------------------------------------

The `require` package provides same global functions as the `assert` package, but instead of returning a boolean result
they terminate current test.

See [t.FailNow](http://golang.org/pkg/testing/#T.FailNow) for details.

[`mock`](http://godoc.org/github.com/stretchr/testify/mock "API documentation") package
----------------------------------------------------------------------------------------

The `mock` package provides a mechanism for easily writing mock objects that can be used in place of real objects when
writing test code.

An example test function that tests a piece of code that relies on an external object `testObj`, can setup
expectations (testify) and assert that they indeed happened:

```go
package yours

import (
	"testing"
	"github.com/stretchr/testify/mock"
)

/*
  Test objects
*/

// MyMockedObject is a mocked object that implements an interface
// that describes an object that the code I am testing relies on.
type MyMockedObject struct {
	mock.Mock
}

// DoSomething is a method on MyMockedObject that implements some interface
// and just records the activity, and returns what the Mock object tells it to.
//
// In the real object, this method would do something useful, but since this
// is a mocked object - we're just going to stub it out.
//
// NOTE: This method is not being tested here, code that uses this object is.
func (m *MyMockedObject) DoSomething(number int) (bool, error) {

	args := m.Called(number)
	return args.Bool(0), args.Error(1)

}

/*
  Actual test functions
*/

// TestSomething is an example of how to use our test object to
// make assertions about some target code we are testing.
func TestSomething(t *testing.T) {

	// create an instance of our test object
	testObj := new(MyMockedObject)

	// setup expectations
	testObj.On("DoSomething", 123).Return(true, nil)

	// call the code we are testing
	targetFuncThatDoesSomethingWithObj(testObj)

	// assert that the expectations were met
	testObj.AssertExpectations(t)

}

// TestSomethingWithPlaceholder is a second example of how to use our test object to
// make assertions about some target code we are testing.
// This time using a placeholder. Placeholders might be used when the
// data being passed in is normally dynamically generated and cannot be
// predicted beforehand (eg. containing hashes that are time sensitive)
func TestSomethingWithPlaceholder(t *testing.T) {

	// create an instance of our test object
	testObj := new(MyMockedObject)

	// setup expectations with a placeholder in the argument list
	testObj.On("DoSomething", mock.Anything).Return(true, nil)

	// call the code we are testing
	targetFuncThatDoesSomethingWithObj(testObj)

	// assert that the expectations were met
	testObj.AssertExpectations(t)

}

// TestSomethingElse2 is a third example that shows how you can use
// the Unset method to cleanup handlers and then add new ones.
func TestSomethingElse2(t *testing.T) {

	// create an instance of our test object
	testObj := new(MyMockedObject)

	// setup expectations with a placeholder in the argument list
	mockCall := testObj.On("DoSomething", mock.Anything).Return(true, nil)

	// call the code we are testing
	targetFuncThatDoesSomethingWithObj(testObj)

	// assert that the expectations were met
	testObj.AssertExpectations(t)

	// remove the handler now so we can add another one that takes precedence
	mockCall.Unset()

	// return false now instead of true
	testObj.On("DoSomething", mock.Anything).Return(false, nil)

	testObj.AssertExpectations(t)
}
```

For more information on how to write mock code, check out
the [API documentation for the `mock` package](http://godoc.org/github.com/stretchr/testify/mock).

You can use the [mockery tool](http://github.com/vektra/mockery) to autogenerate the mock code against an interface as
well, making using mocks much quicker.

[`suite`](http://godoc.org/github.com/stretchr/testify/suite "API documentation") package
-----------------------------------------------------------------------------------------

The `suite` package provides functionality that you might be used to from more common object oriented languages. With
it, you can build a testing suite as a struct, build setup/teardown methods and testing methods on your struct, and run
them with 'go test' as per normal.

An example suite is shown below:

```go
// Basic imports
import (
"testing"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including a T() method which
// returns the current testing context
type ExampleTestSuite struct {
suite.Suite
VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ExampleTestSuite) SetupTest() {
suite.VariableThatShouldStartAtFive = 5
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ExampleTestSuite) TestExample() {
assert.Equal(suite.T(), 5, suite.VariableThatShouldStartAtFive)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
suite.Run(t, new(ExampleTestSuite))
}
```

For a more complete example, using all of the functionality provided by the suite package, look at
our [example testing suite](https://github.com/stretchr/testify/blob/master/suite/suite_test.go)

For more information on writing suites, check out
the [API documentation for the `suite` package](http://godoc.org/github.com/stretchr/testify/suite).

`Suite` object has assertion methods:

```go
// Basic imports
import (
"testing"
"github.com/stretchr/testify/suite"
)

// Define the suite, and absorb the built-in basic suite
// functionality from testify - including assertion methods.
type ExampleTestSuite struct {
suite.Suite
VariableThatShouldStartAtFive int
}

// Make sure that VariableThatShouldStartAtFive is set to five
// before each test
func (suite *ExampleTestSuite) SetupTest() {
suite.VariableThatShouldStartAtFive = 5
}

// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *ExampleTestSuite) TestExample() {
suite.Equal(suite.VariableThatShouldStartAtFive, 5)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestExampleTestSuite(t *testing.T) {
suite.Run(t, new(ExampleTestSuite))
}
```

------

Installation
============

To install Testify, use `go get`:

    go get github.com/stretchr/testify

This will then make the following packages available to you:

    github.com/stretchr/testify/assert
    github.com/stretchr/testify/require
    github.com/stretchr/testify/mock
    github.com/stretchr/testify/suite
    github.com/stretchr/testify/http (deprecated)

Import the `testify/assert` package into your code using this template:

```go
package yours

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSomething(t *testing.T) {

	assert.True(t, true, "True is true!")

}
```

------

Staying up to date
==================

To update Testify to the latest version, use `go get -u github.com/stretchr/testify`.

------

Supported go versions
==================

We currently support the most recent major Go versions from 1.19 onward.

------

Contributing
============

Please feel free to submit issues, fork the repository and send pull requests!

When submitting an issue, we ask that you please include a complete test function that demonstrates the issue. Extra
credit for those using Testify to write the test code that demonstrates it.

Code generation is used. Look for `CODE GENERATED AUTOMATICALLY` at the top of some files. Run `go generate ./...` to
update generated files.

We also chat on the [Gophers Slack](https://gophers.slack.com) group in the `#testify` and `#testify-dev` channels.

------

License
=======

This project is licensed under the terms of the MIT license.
