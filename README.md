# dconf-filter

## Usage

Install:
```
go get github.com/tseho/dconf-filter
```

see "Rules" for the rules file

```
dconf dump / | dconf-filter --rules=<RULES_PATH>
```

## Rules

You need to provide a rules file using the gitignore syntax.

Eg:
```
org/gnome/shell
!org/gnome/shell/command-history
```
This exemple will whitelist all settings in `org/gnome/shell/*` except `org/gnome/shell/command-history`.

As git does with gitignore, __the last matching rule will overwrite__ any previous match.
