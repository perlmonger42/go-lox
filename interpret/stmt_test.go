package interpret

import (
	"github.com/perlmonger42/go-lox/config"
	"github.com/perlmonger42/go-lox/lox"
	"github.com/perlmonger42/go-lox/parse"
	"github.com/perlmonger42/go-lox/resolve"
	"github.com/perlmonger42/go-lox/scan"
)

func exec(text string) {
	config := config.New()
	lox := lox.New(config)
	scanner := scan.New(lox, text)
	tokens := scanner.ScanTokens()
	if lox.HadError {
		return
	}
	parser := parse.New(lox, tokens)
	stmts := parser.ParseProg()
	if lox.HadError {
		return
	}
	interpreter := New(lox)
	var resolver *resolve.T = resolve.New(lox, interpreter)
	resolver.ResolveStmtList(stmts)
	if lox.HadError {
		return
	}

	interpreter.InterpretStmts(stmts)
}

func ExampleEmptyExec() {
	exec("")
	// Output:
}

func ExampleExpressionZeroExec() {
	exec("0;")
	// Output:
}

func ExamplePrintDouglasAdams42() {
	exec("print 6 * 9;")
	// Output:
	// 54
}

func ExamplePrints() {
	exec(`print "one";
          print true;
		  print 2 + 1;
		  print nil;`)
	// Output:
	// one
	// true
	// 3
	// nil
}

func ExampleVar() {
	exec(`var a = 1;
          var b = 2;
		  print a + b;
		  a = b = a + b;
		  print a + b;
		  c = 7;`)
	// Output:
	// 3
	// 6
	// [line 6] Error at 'Identifier': Undefined variable 'c'.
}

func ExampleBlocks() {
	exec(`var a = "global a";
          var b = "global b";
          var c = "global c";
          {
            var a = "outer a";
            var b = "outer b";
            {
              var a = "inner a";
              print a;
              print b;
              print c;
            }
            print a;
            print b;
            print c;
          }
          print a;
          print b;
          print c;`)
	// Output:
	// inner a
	// outer b
	// global c
	// outer a
	// outer b
	// global c
	// global a
	// global b
	// global c
}

func ExampleRedefine() {
	exec(`var a = 1;
          print a;
          var a = 2; // redefining variables is allowed in the global scope
          print a;
          {
            var a = 3;
			print a;
			var a = 4; // but not in nested scopes
			print a;
		  }`)
	// Output:
	// 1
	// 2
	// 3
	// [line 8] Error at 'Identifier': Variable 'a' redefined.
	// 4
}

func ExampleUndefinedVar() {
	exec(`var x; print x;`)
	// Output:
	// nil
}

func ExampleTheVarIsDefinedAfterItsInitializerIsEvaluated() {
	exec(`var a = 1;
          {
            var a = a + 2;
            print a;
          }`)
	// Output:
	// 3
}

func ExampleIfStatement() {
	exec(`var x = 3;
	      if (x < 3) print "smaller than 3";
		  else if (x > 3) print "larger than 3";
		  else print "exactly 3";
		  `)
	// Output:
	// exactly 3
}

func ExampleAscii() {
	exec(`var Testing =     1;
          var testing =    20;
		  var one     =   300;
		  var two     =  4000;
		  var three   = 50000;
	      print Testing + testing + one + two + three;
	`)
	// Output:
	// 54321
}

func ExampleUnicode() {
	// "Ťėšťǐňġ" is entirely made of 2-byte UTF-8 encodings
	// "ṫẹṡṫịṅḡ" is entirely made of 3-byte UTF-8 encodings
	// "𝕠𝕟𝕖, 𝕥𝕨𝕠, 𝕥𝕙𝕣𝕖𝕖" has a bunch of 4-byte UTF-8 encodings
	exec(`var Ťėšťǐňġ =     1;
          var ṫẹṡṫịṅḡ =    20;
		  var 𝕠𝕟𝕖     =   300;
		  var 𝕥𝕨𝕠     =  4000;
		  var 𝕥𝕙𝕣𝕖𝕖   = 50000;
	      print Ťėšťǐňġ + ṫẹṡṫịṅḡ + 𝕠𝕟𝕖 + 𝕥𝕨𝕠 + 𝕥𝕙𝕣𝕖𝕖;
	`)
	// Output:
	// 54321
}

func ExampleWhileWithZeroIterations() {
	exec(`var n = 10;
          while (n < 10) print n;
		  print "done";`)
	// Output:
	// done
}

func ExampleWhileWithOneIteration() {
	exec(`var n = 9;
          while (n < 10) n = n + 1;
		  print n;
		  print "done";`)
	// Output:
	// 10
	// done
}

func ExampleWhileWithTwoIterations() {
	exec(`var n = 8;
          while (n < 10) n = n + 1;
		  print n;
		  print "done";`)
	// Output:
	// 10
	// done
}

func ExampleWhileWithZeroIterationsUsingBlock() {
	exec(`var n = 10;
          while (n < 10) { print n; n = n + 1; }
		  print "done";`)
	// Output:
	// done
}

func ExampleWhileWithOneIterationUsingBlock() {
	exec(`var n = 0;
          while (n < 1) { print n; n = n + 1; }
		  print "done";`)
	// Output:
	// 0
	// done
}

func ExampleWhileWithTwoIterationsUsingBlock() {
	exec(`var n = 0;
          while (n < 2) { print n; n = n + 1; }
		  print "done";`)
	// Output:
	// 0
	// 1
	// done
}

func ExampleWhileCountdown() {
	exec(`var n = 10;
          while (n > 0) {
		    print n;
			n = n - 1 ;
	      }
		  print "liftoff!";`)
	// Output:
	// 10
	// 9
	// 8
	// 7
	// 6
	// 5
	// 4
	// 3
	// 2
	// 1
	// liftoff!
}

func ExampleNestedAndUnnest() {
	exec(`var a = "global a";
          {
            var b = "outer b";
            {
              var c = "inner c";
              print a + ", " + b + ", " + c;
            }
            print a + ", " + b + ", " + c;
		  }
          print a + ", " + b + ", " + c;
		  `)
	// Output:
	// global a, outer b, inner c
	// [line 8] Error at 'Identifier': Undefined variable 'c'.
	// global a, outer b, {([<nil>])}
	// [line 10] Error at 'Identifier': Undefined variable 'b'.
	// [line 10] Error at 'Identifier': Undefined variable 'c'.
	// global a, {([<nil>])}, {([<nil>])}
}

func ExampleNestedUndefined() {
	exec(`var a = "global a";
          {
            var b = "outer b";
            {
              var c = "inner c";
              print a;
              print b;
              print c;
			  print d;
            }
		  }`)
	// Output:
	// global a
	// outer b
	// inner c
	// [line 9] Error at 'Identifier': Undefined variable 'd'.
	// nil
}

func ExampleFibonacciWhileLoop() {
	exec(`var a = 0; var b = 1;
          while (a < 10000) {
		    print a;
			var temp = a;
			a = b;
			b = temp + b;
          }`)
	// Output:
	// 0
	// 1
	// 1
	// 2
	// 3
	// 5
	// 8
	// 13
	// 21
	// 34
	// 55
	// 89
	// 144
	// 233
	// 377
	// 610
	// 987
	// 1597
	// 2584
	// 4181
	// 6765
}

func ExampleForLoop() {
	exec(`for (var q = 1; q < 10000; q = q * 2) print q;`)
	// Output:
	// 1
	// 2
	// 4
	// 8
	// 16
	// 32
	// 64
	// 128
	// 256
	// 512
	// 1024
	// 2048
	// 4096
	// 8192
}

func ExampleBadReturn() {
	exec(`return 17;`)
	// Output:
	// [line 1] Error at 'Return': Cannot return from top-level code.
}

func ExampleFunction() {
	exec(`fun sayHi(first, last) {
            print "Hi, " + first + " " + last + "!";
          }
		  sayHi("Dear", "Reader");`)
	// Output:
	// Hi, Dear Reader!
}

func ExampleSquare() {
	exec(`fun square(n) {
            return n * n;
          }
		  print square(7);`)
	// Output:
	// 49
}

func ExampleFibonacciRecursive() {
	exec(`fun fib(n) { if (n <= 1) return n; return fib(n - 2) + fib(n - 1); } for (var i = 0; i < 20; i = i + 1) { print fib(i); }`)
	// Output:
	// 0
	// 1
	// 1
	// 2
	// 3
	// 5
	// 8
	// 13
	// 21
	// 34
	// 55
	// 89
	// 144
	// 233
	// 377
	// 610
	// 987
	// 1597
	// 2584
	// 4181
}

func ExampleClosure() {
	exec(`
		fun makeCounter() {
		  var i = 0;
		  fun count() {
			i = i + 1;
			print i;
		  }

		  return count;
		}
		var counter = makeCounter();
		counter(); // "1".
		counter(); // "2"
	`)
	// Output:
	// 1
	// 2
}

func ExampleRedeclareFormal() {
	exec(`
		fun f(a) {
		  print a;
		  var a = "redefined";
		  print a;
		}
		f(13);
	`)
	// Output:
	// 13
	// [line 4] Error at 'Identifier': Variable 'a' redefined.
	// redefined
}

func ExampleVariableResolution() {
	exec(`
var a = "global";
{
  fun showA() {
    print a;
  }

  showA();
  var a = "local";
  showA();
}
	`)
	// Output:
	// global
	// global
}

func ExampleVariableInitCannotReferenceItself() {
	exec(`
var a = "outer";
{
  var a = "[" + a + "]";
  print a;
}
print a;
	`)
	// Originally, output was:
	//   [line 4] Error at 'Identifier': Cannot read local variable in its own initializer.
	// but now,
	// Output:
	// [outer]
	// outer
}

func ExampleInvalidThis() {
	exec(`
print this;
	`)
	// Output:
	// [line 2] Error at 'This': Cannot use `this` outside of a class.
}

func ExampleClass() {
	exec(`
class DevonshireCream {
  serveOn() {
    return "Scones";
  }
}
print DevonshireCream;
	`)
	// Output:
	// class DevonshireCream
}

func ExampleInstance() {
	exec(`
class Bagel { }
print Bagel();
	`)
	// Output:
	// Bagel{}
}

func ExampleSetGet() {
	exec(`
class Widget { }
var w = Widget();
w.name = "my shiny new widget";
w.weight = 42;
print w.name;
print w.weight;
print w.nonexistent;
	`)
	// Output:
	// my shiny new widget
	// 42
	// [line 8] Error at 'Identifier': Undefined property `nonexistent`.
	// runtime error: {Identifier: `nonexistent` Undefined property `nonexistent`.}
}

func ExampleMethod() {
	exec(`
class Bacon {
  eat() {
    print "Crunch crunch crunch";
  }
}
Bacon().eat();
	`)
	// Output:
	// Crunch crunch crunch
}

func ExampleCake() {
	exec(`
class Cake {
  taste() {
    var adjective = "delicious";
    print "The " + this.flavor + " cake is " + adjective + "!";
  }
}
var cake = Cake();
cake.flavor = "German chocolate";
cake.taste();
	`)
	// Output:
	// The German chocolate cake is delicious!
}

func ExampleMethodPointer() {
	exec(`
class Cake {
  taste() {
    var adjective = "delicious";
    print "The " + this.flavor + " cake is " + adjective + "!";
  }
}
var cake = Cake();
cake.flavor = "Devil's Food";
var f = cake.taste;
f();
	`)
	// Output:
	// The Devil's Food cake is delicious!
}

func ExampleConstructor() {
	exec(`
class Point {
  init(x, y) {
    this.x = x;
	this.y = y;
  }
  str() {
    return "(" + str(this.x) + "," + str(this.y) + ")";
  }
}
var pyth = Point(3, 4);
print pyth.str();
	`)
	// Output:
	// (3,4)
}

func ExampleIllegalReturnValueFromInit() {
	exec(`
class Datum {
  init() {
    return 7;
  }
}
	`)
	// Output:
	// [line 4] Error at 'Return': Cannot return a value from an initializer.
}

func ExampleReturnEarlyFromInit() {
	exec(`
class Datum {
  init(x) {
    this.number = x;
    if (x == 0) {
	  this.name = "zero";
	  return;
	}
	this.name = "nonzero";
  }
  str() {
    return this.name + "(" + str(this.number) + ")";
  }
}
print Datum(0).str();
print Datum(1).str();
	`)
	// Output:
	// zero(0)
	// nonzero(1)
}

func ExampleSuperclass() {
	exec(`
class Base {}
class Extension < Base {}
print Base;
print Extension;
	`)
	// Output:
	// class Base
	// class Extension
}

func ExampleCyclicInheritance() {
	exec(`
class C < C {
}
	`)
	// Output:
	// [line 2] Error at 'Identifier': A class can't inherit from itself.
}

func ExampleNonclassSuperclass() {
	exec(`
var NotAClass = "I am totally not a class.";
class Extension < NotAClass { }
	`)
	// Output:
	// [line 3] Error at 'Identifier': Superclass must be a class.
	// runtime error: {Identifier: `NotAClass` Superclass must be a class.}
}

func ExampleUndefinedSuperclass() {
	exec(`
class Extension < Base {
}
	`)
	// Output:
	// [line 2] Error at 'Identifier': Undefined variable 'Base'.
	// [line 2] Error at 'Identifier': Superclass must be a class.
	// runtime error: {Identifier: `Base` Superclass must be a class.}
}

func ExampleInheritMethod() {
	exec(`
class Doughnut {
  cook() {
    print "Fry until golden brown.";
  }
}
class BostonCream < Doughnut {
  cook() {
    super.cook();
	print "Pipe full of custard and coat with chocolate.";
  }
}
BostonCream().cook();
	`)
	// Output:
	// Fry until golden brown.
	// Pipe full of custard and coat with chocolate.
}

func ExampleBoundFunction() {
	exec(`
class Doughnut {
  cook() {
    print "Fry until golden brown.";
  }
}
var cookIt = Doughnut().cook;
cookIt();
	`)
	// Output:
	// Fry until golden brown.
}

func ExampleSuperIsRelativeToMethodNotInstance() {
	exec(`
class A {
  method() {
    print "Method 'A'.";
  }
}

class B < A {
  method() {
    print "Method 'B'.";
  }

  test() {
    super.method();
  }
}

class C < B {}

C().test();
	`)
	// Output:
	// Method 'A'.
}

func ExampleBoundSuper() {
	exec(`
class Doughnut {
  cook() {
    print "Fry until golden brown.";
  }
}
class BostonCream < Doughnut {
  cook() {
    var cooker = super.cook;
	cooker();
	print "Pipe full of custard and coat with chocolate.";
  }
}
BostonCream().cook();
	`)
	// Output:
	// Fry until golden brown.
	// Pipe full of custard and coat with chocolate.
}

func ExampleInvalidUsesOfSuper() {
	exec(`
class Eclair {
  cook() {
    super.cook();
    print "Pipe full of crème pâtissière.";
  }
}
super.notEvenInAClass();
	`)
	// Output:
	// [line 4] Error at 'Super': Can't use 'super' in a class with no superclass.
	// [line 8] Error at 'Super': Can't use 'super' outside of a class.
}
