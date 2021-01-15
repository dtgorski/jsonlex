#### v0.3.1
* Test and documentation improvements 

#### v0.3.0
* Lexer.Scan() is now reentrant. Depending on the outcome of the Yield callback, the Scan() function will terminate and can be invoked again. 
* Added a Cursor which provides iteration over tokens using Curr(), Peek(), Next() and Last() methods
* Added a Filter, an optional callback for the Cursor constructor, which allows to filter the lexer output and transparently skip tokens.
* API break:
  * Old: Yield func(token Token, load []byte, pos uint)
  * New: Yield func(kind TokenKind, load []byte, pos uint) bool

#### v0.2.4
* First milestone
