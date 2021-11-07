# Things Catalog Analyzer

I recently started catalogging all the things I own. I simply wanted to have
an overview over my things, think about what I can get rid off and maybe even
plot the data onto a graph over a longer timer to see how many things i throw
away or buy.

## File Format

The format is pretty simple, here's an example:

```
## Category 1

Item A,5
Item B,1
Item C,3

## Category

Item D,4
```

The number after each item is the count.

## Usage

In a terminal, run:

```
things-catalog-analyzer FILE_PATH
```

You should then see something like this:

```
Count    Name
----     ----
42       Hygiene
11       Entertainment
15       Möbel
7        Putzutensil
8        Werkzeug
8        Gesundheit
13       Schlaf
3        Komfort
38       Tech Entertainment
13       Misc
89       Kleidung
30       Küchenutensil
90       Geschirr
7        Dekoration


Overall Item count: 374
```
