// Released under an MIT-style license. See LICENSE.

// Generated by generate.oh

package boot

var Script string = `
define _connect_: syntax (conduit name) = {
	set conduit: eval conduit
	syntax (left right) e = {
		define p: conduit
		spawn {
			e::eval: quasiquote: block {
				public (unquote name) = (unquote p)
				eval (unquote left)
			}
			p::_writer_close_
		}
		block {
			define s: e::eval: quasiquote: block {
				public _stdin_ = (unquote p)
				eval (unquote right)
			}
			p::_reader_close_
			return s
		}
	}
}
define _redirect_: syntax (name mode closer) = {
	syntax (c cmd) e = {
		define c: e::eval c
		define f = ()
		if (not: or (is-channel c) (is-pipe c)) {
			set f: open mode c
			set c = f
		}
		define s: e::eval: quasiquote: block {
			public (unquote name) (unquote c)
			eval (unquote cmd)
		}
		if (not: is-null f): eval: quasiquote: f::(unquote closer)
		return s
	}
}
define ...: method (: args) = {
	cd _origin_
	define path: args::head
	if (eq 2: args::length) {
		cd path
		set path: args::get 1
	}
	while true {
		define abs: symbol: "/"::join $PWD path
		if (exists abs): return abs
		if (eq $PWD /): return path
		cd ..
	}
}
define and: syntax (: lst) e = {
	define r = false
	while (not: is-null lst) {
		set r: e::eval: lst::head
		if (not r): return r
		set lst: lst::tail ()
	}
	return r
}
define _append_stderr_: _redirect_ _stderr_ "a" _writer_close_
define _append_stdout_: _redirect_ _stdout_ "a" _writer_close_
define _backtick_: syntax (cmd) e = {
	define p: pipe
	spawn {
		e::eval: quasiquote: block {
			public _stdout_ = (unquote p)
			eval (unquote cmd)
		}
		p::_writer_close_
	}
	define r: cons () ()
	define c = r
	while (define l: p::readline) {
		set-cdr c: cons l ()
		set c: cdr c
	}
	p::_reader_close_
	return: cdr r
}
define catch: syntax (name: clause) e = {
	define args: list name (quote throw)
	define body: list (quote throw) name
	if (is-null clause) {
		set body: list body
	} else {
		set body: append clause body
	}
	define handler: e::eval {
		list (quote method) args (quote =) @body
	}
	define _return: e::eval (quote return)
	define _throw = throw
	e::public throw: method (condition) = {
		_return: handler condition _throw
	}
}
define _channel_stderr_: _connect_ channel _stderr_
define _channel_stdout_: _connect_ channel _stdout_
define coalesce: syntax (: lst) e = {
	while (and (not: is-null: cdr lst) (not: resolves: lst::head)) {
		set lst: cdr lst
	}
	return: e::eval: lst::head
}
define echo: builtin (: args) = {
	if (is-null args) {
		_stdout_::write: symbol ""
	} else {
		_stdout_::write @(for args symbol)
	}
}
define error: builtin (: args) =: _stderr_::write @args
define for: method (l m) = {
	define r: cons () ()
	define c = r
	while (not: is-null l) {
		set-cdr c: cons (m: l::head) ()
		set c: cdr c
		set l: cdr l
	}
	return: cdr r
}
define glob: builtin (: args) =: return args
define import: syntax (name) e = {
	set name: e::eval name
	define m: module name
	if (or (is-null m) (is-object m)) {
		return m
	}

	e::eval: quasiquote: _root_::define (unquote m): object {
		source (unquote name)
	}
}
define is-list: method (l) = {
	if (is-null l): return false
	if (not: is-cons l): return false
	if (is-null: cdr l): return true
	is-list: cdr l
}
define is-text: method (t) =: or (is-string t) (is-symbol t)
# TODO: Replace with builtin rather than invoking bc.
define math: method (S) e = {
	set S: e::interpolate S
	catch ex {
		set ex::type = "error/syntax"
		set ex::message = $"Malformed expression: ${S}"
	}

	float @(_backtick_ (block {
		echo "scale=6"
		write: symbol S
	} | bc))
}
define object: syntax (: body) e = {
	e::eval: cons (quote block): append body (quote: context)
}
define or: syntax (: lst) e = {
	define r = false
	while (not: is-null lst) {
		set r: e::eval: lst::head
		if r: return r
		set lst: lst::tail ()
	}
	return r
}
define _pipe_stderr_: _connect_ pipe _stderr_
define _pipe_stdout_: _connect_ pipe _stdout_
define printf: method (f: args) =: echo: f::sprintf @args
define _process_substitution_: syntax (:args) e = {
	define fifos = ()
	define procs = ()
	define cmd: for args: method (arg) = {
		if (not: is-cons arg): return arg
		if (eq (quote _substitute_stdin_) (arg::head)) {
			define fifo: temp-fifo
			define proc: spawn {
				e::eval: quasiquote {
					_redirect_stdin_ {
						unquote fifo
						unquote: cdr arg
					}
				}
			}
			set fifos: cons fifo fifos
			set procs: cons proc procs
			return fifo
		}
		if (eq (quote _substitute_stdout_) (arg::head)) {
			define fifo: temp-fifo
			define proc: spawn {
				e::eval: quasiquote {
					_redirect_stdout_ {
						unquote fifo
						unquote: cdr arg
					}
				}
			}
			set fifos: cons fifo fifos
			set procs: cons proc procs
			return fifo
		}
		return arg
	}
	e::eval cmd
	wait @procs
	rm @fifos
}
define quasiquote: syntax (cell) e = {
	if (not: is-cons cell): return cell
	if (is-null cell): return cell
	if (eq (quote unquote): cell::head): return: e::eval: cell::get 1
	cons {
		e::eval: list (quote quasiquote): cell::head
		e::eval: list (quote quasiquote): cdr cell
	}
}
define quote: syntax (cell) =: return cell
define read: builtin () =: _stdin_::read
define readline: builtin () =: _stdin_::readline
define _redirect_stderr_: _redirect_ _stderr_ "w" _writer_close_
define _redirect_stdin_: _redirect_ _stdin_ "r" _reader_close_
define _redirect_stdout_: _redirect_ _stdout_ "w" _writer_close_
define source: syntax (name) e = {
	define basename: e::eval name
	define paths = ()
	define name = basename

	if (has $OHPATH): set paths: ":"::split $OHPATH
	while (and (not: is-null paths) (not: exists name)) {
		set name: "/"::join (paths::head) basename
		set paths: cdr paths
	}

	if (not: exists name): set name = basename

	define f: open r- name

	define r: cons () ()
	define c = r
	while (define l: f::read) {
		set-cdr c: cons (cons (get-line-number) l) ()
		set c: cdr c
	}
	set c: cdr r
	f::close

	define rval: status 0
	define eval-list: method (first rest) o = {
		if (is-null first): return rval
		set-line-number: first::head
		set rval: e::eval: cdr first
		eval-list (rest::head ()) (cdr rest)
	}
	eval-list (c::head ()) (cdr c)
	return rval
}
define write: method (: args) =: _stdout_::write @args
_sys_::public exception: method (type message status file line) = {
	object {
		public type = type
		public status = status
		public message = message
		public line = line
		public file = file
	}
}
_sys_::public get-prompt: method self (suffix) = {
	catch unused {
		return suffix
	}
	self::prompt suffix
}
_sys_::public prompt: method (suffix) = {
	define dirs: "/"::split $PWD
	return: ""::join (dirs::get -1) suffix
}
_sys_::public throw: method (c) = {
	error: ": "::join c::file c::line c::type c::message
	fatal c::status
}

exists ("/"::join $HOME .ohrc) && source ("/"::join $HOME .ohrc)

`

//go:generate ./generate.oh
//go:generate go fmt generated.go
