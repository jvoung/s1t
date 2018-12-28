# S1t

## Overview

A simple SAT solver, written for fun (learn more Go lang).

## Input

Initially, this will handle input text in the format of DIMACS-CNF:
<http://www.domagoj-babic.com/uploads/ResearchProjects/Spear/dimacs-cnf.pdf>

May extend this to 7-bit ascii CNF format like:
`(x1 | ~x5 | x2) & (~x1 | x5 | x3 | x4)`
Over time, may extend this to non-CNF SAT formats.

## Naming (or, why s1t?)

s1t is a silly and trivial abbreviation in the style of i18n, l10n, S12n.
You can guess the missing letter.
