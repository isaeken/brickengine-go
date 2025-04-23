# 📚 BrickEngine Built-in Function Reference

This document outlines all available built-in functions grouped by category, along with usage examples.

---

## 🆔 UUID & Formatting

| Function           | Description                        | Example                             |
|--------------------|------------------------------------|-------------------------------------|
| `uuid()`           | Generates a new UUID               | `let id = uuid()`                   |
| `slug(value)`      | Converts string to slug            | `slug("My Title") → "my-title"`     |
| `random()`         | Returns a random number [0.0–1.0)  | `random() → 0.427`                  |
| `now()`            | Current date/time string           | `now() → "2024-04-22T19:00:00Z"`    |
| `format(value)`    | Converts value to string           | `format({ a: 1 }) → '{"a":1}'`      |
| `to_json(value)`   | JSON stringifies a value           | `to_json([1,2]) → "[1,2]"`          |
| `parse_json(str)`  | Parses a JSON string               | `parse_json('{"a":1}').a → 1`       |

---

## 🔠 String Functions

| Function                   | Description                            | Example                                |
|----------------------------|----------------------------------------|----------------------------------------|
| `strlen(str)`              | Length of the string                   | `strlen("isa") → 3`                    |
| `str_upper(str)`           | Converts to uppercase                  | `str_upper("isa") → "ISA"`             |
| `str_lower(str)`           | Converts to lowercase                  | `str_lower("ISA") → "isa"`             |
| `str_trim(str)`            | Trims whitespace                       | `str_trim("  isa  ") → "isa"`          |
| `str_contains(str, sub)`   | Checks substring                       | `str_contains("isa", "a") → true`      |
| `str_starts_with(str, x)`  | Starts with check                      | `str_starts_with("isa", "i") → true`   |
| `str_ends_with(str, x)`    | Ends with check                        | `str_ends_with("isa", "a") → true`     |
| `str_replace(str, x, y)`   | Replaces x with y                      | `str_replace("a-b", "-", "_") → "a_b"` |
| `substr(str, start, len)`  | Substring                              | `substr("isaeken", 0, 3) → "isa"`      |
| `split(str, sep)`          | Splits into array                      | `split("a,b,c", ",") → ["a","b","c"]`  |
| `join(arr, sep)`           | Joins array into string                | `join(["a","b"], "-") → "a-b"`         |
| `repeat(str, n)`           | Repeats string n times                 | `repeat("a", 3) → "aaa"`               |
| `str_reverse(str)`         | Reverses the string                    | `str_reverse("abc") → "cba"`           |

---

## ➗ Math Functions

| Function         | Description                            | Example               |
|------------------|----------------------------------------|------------------------|
| `abs(x)`         | Absolute value                         | `abs(-5) → 5`          |
| `round(x)`       | Rounds to nearest integer              | `round(2.7) → 3`       |
| `floor(x)`       | Rounds down                            | `floor(2.9) → 2`       |
| `ceil(x)`        | Rounds up                              | `ceil(2.1) → 3`        |
| `min(a, b)`      | Minimum of two numbers                 | `min(3, 9) → 3`        |
| `max(a, b)`      | Maximum of two numbers                 | `max(3, 9) → 9`        |
| `sqrt(x)`        | Square root                            | `sqrt(9) → 3`          |
| `pow(a, b)`      | Power                                  | `pow(2, 3) → 8`        |

---

## 🧠 Type & Meta

| Function        | Description                 | Example                           |
|------------------|-----------------------------|------------------------------------|
| `type_of(value)` | Returns type as a string    | `type_of([1, 2]) → "array"`        |

---

## 📦 Array Functions

| Function               | Description                            | Example                                |
|------------------------|----------------------------------------|----------------------------------------|
| `count(array)`         | Returns number of elements             | `count([1,2,3]) → 3`                    |
| `push(array, value)`   | Adds value to end                      | `push([1], 2) → [1,2]`                  |
| `pop(array)`           | Removes last element                   | `pop([1,2]) → [1]`                      |
| `shift(array)`         | Removes first element                  | `shift([1,2]) → [2]`                    |
| `unshift(array, val)`  | Adds value to beginning                | `unshift([2], 1) → [1,2]`               |
| `includes(array, val)` | Checks if value is in array            | `includes([1,2], 2) → true`             |
| `index_of(array, val)` | Index of value                         | `index_of(["a","b"], "b") → 1`          |
| `reverse(array)`       | Reverses array                         | `reverse([1,2]) → [2,1]`                |
| `sort(array)`          | Sorts numbers (ascending)              | `sort([3,1,2]) → [1,2,3]`               |
| `slice(arr, start, end)`| Slices array (exclusive end)         | `slice([0,1,2,3],1,3) → [1,2]`          |
| `concat(arr1, arr2)`   | Merges two arrays                      | `concat([1], [2,3]) → [1,2,3]`          |

---

> 💡 This list grows as BrickEngine evolves. You can register your own native functions via Go runtime too.
