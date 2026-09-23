# PNG fixtures

`generate.py` uses only Python's `struct`, `zlib` and `pathlib` modules to
assemble deterministic PNG files. It computes independent zlib streams and
PNG CRC32 checksums, packs sub-byte samples, applies all five PNG filters,
and emits each Adam7 pass directly. Regenerate with
`python3 png/tests/data/generate.py` from the repository root.
The malformed files retain valid CRCs unless their names explicitly test CRC
handling, so tests reach the intended parser branch.
