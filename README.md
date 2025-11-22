# GoReactiveHTML (aka GoReacTML)
Like we have pythonic way, we are creating a Gothonic way for webdevelopment.

GoReactiveHTML is a HTML on Gothonic way. 

The best of Javascript, css and HTML transparented integrated on Go, with less code and more reactive results.

Focus on GoReactiveHTML code and Go. And let's make web development extremly productive and simple again.

## Goal
Project for an extended version of HTML with native Go Language integration. HTML in a GOThonic way

### What is it look like
Go is the simple, fast e light.
HTML, JS and CSS are the web"languages"
To code Frontend + Backend separated causes a lot of problems, like:
- Code/Login repetition
- Error-prone integrations
- Many Unproductive Request
- Many Data Conversion Name or structure

Our propose is to avoid repetition and create a transparent reactive integration to code less and produce more results for final user in solutions based on HTML. Now your HTML will be alive and integrated to powerful backend powered by Go Language

## Examples
SHOW ME THE CODE
- Todo app
- Count app
- Sign in page
- Sign up page
- Home Broker app
- Contact Form
- Celsius to Fahrenheit page
- Upload file


### TDLR; 

To make web development extremely productive and simple again using Go and GoReactiveHTML, here are key points and approaches based on current tools and practices:

#### Go for Web Development
* Performance & Simplicity: Go is a statically typed, compiled language that offers fast execution and a clean, easy-to-learn syntax, reducing the learning curve for web development

* Concurrency: Built-in concurrency with goroutines and channels allows scalable and high-performing web apps.

* Standard Library: Go’s extensive standard library supports building web servers and handling HTTP requests without needing many third-party dependencies

* Cross-Platform: Compiled Go binaries run across different operating systems without modification

#### GoReactiveHTML for Reactive Web Components

What is GoReactiveHTML? It is a way to build custom reactive HTML components integrated to Go applications, allowing you to interact HTML, CSS, and JavaScript templates transparently Front and back directly in HTML or JS code with dynamic bindings

.

#### How It Works

At the end, when it all compiles, we will have just Vanilla Javascript on the front. 

No Complex JS Frameworks or Tools imported: This approach avoids heavy JavaScript frameworks or build steps, making development simpler and more maintainable

Components can dynamically update its title, contents, state from Go. And it can update Go state too.

To make web development extremely productive and simple again with Go and GoReactiveHTML:

* Leverage Go’s performance, concurrency, and simplicity for backend logic and HTTP serving.

* Build reactive UI components using GoReactiveHTML by defining HTML templates with dynamic bindings and event handling in Go and JS, avoiding complex JS frameworks.

* This approach streamlines development by using Go backend and JS on frontend, reduces context switching, and keeps the stack simple and maintainable. You can think it as a Go template extension with reactive superpower. HTML interface data oriented.

This combination modernizes Go web development with reactive UI capabilities, making it both powerful and accessible

Check examples to have simples collapsible details pages/components that can be built with minimal code and full integration

### Dev Dependency
 - extHTML
 - Go

# Unit test
This is a complex system having JS and Go language code, as it target to execute browser communication as client with the server.

All the tests should be executed inside the project container. If you are using devcontainer and vscode, you can use the regular vscode terminal connected to DevContainer.

## JS lib Unit Test
Go to root directory of the project
```
cd web/lib
npm test -- --verbose --runInBand
```

## Go language encode Unit Test
Go to root directory of the project
```
cd tests/encoder
go test go_test.go
```