# Simpl

Simpl is a simple programming language created with the purpose of practicing an learning about different concepts of programming. Following the guide of the book "Writing an interpreter in Go", written by Thorsten Ball, the Simpl language was created using Golang.

This language aims to be simple and provide with a set of tools, functions and methods to create scripts and to be more than just a practice project.

This version is just a working language with basic functions, but it still needs more polishing and newer functionalities in the future.

To use it, download the code, create a file with the smp extension withing the same directory and run the following command:

```
go run main.go <FILE_NAME>
```

With the download there is a file called 'example.smp' that has an example of the language, which you can run and modify with the command above, subtituting <FILE_NAME> with 'example.smp' and it will show the output for he code.

## Installation

Right now there is no installation

## Usage

### Declaring variables

```
var int b = 4;
```

### Conditionals

```
var int a = 0;
var int b = 1;

if (a < b) {
  return 10;
} else {
  return 20;
}
```

### Functions and function calls

```
func add(x, y) {
  return x + y;
}

var int c = add(4, 4);
```

```
var int add = func(x, y) {
  return x + y;
}

var int c = add(4, 4);
```

### For loops

```
var int i = 0;
var array myArray = [];

for (i < 10) {
  push(myArray, i);
  i = i + 1;
}
```

## Built in functions

### Print

Receives any data type and prints it to the console. Can receive any number of arguments

print(<arg>)

```
print("Hello world!");
//  outputs "Hello world!"
```

### Length

Receives a string or array data type and returns the length

length(<string | array>)

```
var string myString = "Hello";
var int len = length(myString);
print(len);
//  outputs 5
```

```
var array myArray = [1, 2, 3];
var int len = length(myArray);
print(len);
//  outputs 3
```

### Push

Adds an element to the end of the array

push(<array_var>, <new_element>)

```
var array myArray = [1, 2, 3];
push(myArray, 4);
print(myArray);
//  outputs [1, 2, 3, 4]
```

### RemoveAt

Removes an element at the specified position on the array

removeAt(<array_var>, <index_to_remove>)

```
var array myArray = [1, 2, 3];
removeAt(myArray, 1);
print(myArray);
//  output [1, 3]
```

### RemoveLast

Removes last element of the array

removeLast(<array_var>)

```
var array myArray = [1, 2, 3];
removeLast(myArray);
print(myArray);
//  outputs [1, 2]
```

### FirstElement

Returns the first element of the array

firstElement(<array_var>)

```
var array myArray = ["Welcome", "to", "Simpl"];
var string element = firstElement(myArray);
print(element);
//  outputs Welcome
```

### LastElement

Returns the last element of the array

lastElement(<array_var>)

```
var array myArray = ["Welcome", "to", "Simpl"];
var string name = lastElement(myArray);
print(name);
//  outputs Simpl
```

### Copy

Makes a copy of an array

copy(<array_var>)

```
var array myArray = [1, 2, 3];
var array mySecondArray = copy(myArray);
print(mySecondArray);
//  outputs [1, 2, 3]
```

### Range

Takes one or two numbers and makes an array between those. If only one is provided, the range goes from zero (0) to the end parameter

range(<start>, <end>) | from start to end
range(<end>) | from 0 to end

```
var array myArray = range(5, 10);
print(myArray);
//  outputs [5, 6, 7, 8, 9, 10]
```

```
var array myArray = range(5);
print(myArray);
//  outputs [0, 1, 2, 3, 4, 5]
```

## Contributing

Right now this is not an open source project
