# Mark V. Shaney

English text isn’t complete random obviously. If you see the words, No matter there’s a pretty good chance they will be followed by the word how. In the mid-1980s, Rob Pike and colleagues used this to write a program to create fake-but-believable posts on internet forums. They did this by writing a program that would read in a large amount of English text and for every pair of words they observed, they computed the frequency with which they saw each 3rd word following them. For example, given the text:

> no matter how hard you try no matter can escape a black hole


They would compute the following table of frequencies:


![Imgur](http://i.imgur.com/iBlehw3.png)


The special empty string entries "" indicate places where there was no previous word (i.e. the start of the input).


Once you have such a table, you can generate English-like text in the following way:

* Choose a random pair of words to start with, say x y.
* Repeatedly choose the next word to output by looking at the table for the two previous output words and picking a word in proportion to how frequently you observed that word in that context.
