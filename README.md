# Silence

At first I want to build a native gui client for pocket on linux, I want to use
golang,but I don't found a suitable GUI library,So I turned to write this little
tool.

## overview
![](https://raw.github.com/hangyan/silence/master/images/usage.png)




## install

`go get github.com/hangyan/go-pocket`

## usage

### authentication

The first thing you need to do is authentication,and this only need to do
once.Just use `go-pocket auth` and you will know what to do.

The authentication will use your default browser to work.

### get

![](https://raw.github.com/hangyan/silence/master/images/get-help.png)

`get` use various filters to retrieve items from pocket and output the result,like:



then you can operate on these items by specify the number (eg:1,2...).

### add

![](https://raw.github.com/hangyan/silence/master/images/get.png)


This project is still under development,any feedback will be welcome.
Email : [yanhangyhy@gmail.com](mailto:yanhangyhy@gmail.com)
