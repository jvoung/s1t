Overview
--------

A simple SAT solver, written as an exercise to learn Go lang and possibly other languages. Initially this will start with the basic DPLL algorithm,
and enhance.

Input
-----

Initially, this will handle input text in the format of DIMACS-CNF:
http://www.domagoj-babic.com/uploads/ResearchProjects/Spear/dimacs-cnf.pdf

May extend this to 7-bit ascii CNF format like:
`(x1 | ~x5 | x4) & (~x1 | x5 | x3 | x4)`
Over time, may extend this further to non-CNF SAT formats.

Naming (or, why s1t?)
---------------------

s1t is a silly and trivial abbreviation in the style of i18n, l10n, S12n.
You can guess the missing letter. Sadly, it looks like "sit".
