What if the neosearch commands can be piped between each other?

(USING name.idx GET "neoway") INTERSECTS (USING uf.idx GET "sc")

(define names (index "name.idx"))
(define states (index "state.idx"))

(and (get names "neoway") (get states "sc"))

(put names ('neoway '(1 2 3 4 5)))
OK
(get names "neoway")
'(1 2 3 4 5)
