# Code Languages Showcase

Short, valid snippets in a broad set of highlight.js-supported languages
(https://highlightjs.readthedocs.io/en/latest/supported-languages.html),
each implementing the same idea: check whether a review decision is approved.

## JavaScript

```javascript
function isApproved(review) {
  return review.status === 'approved'
}
```

## TypeScript

```typescript
interface Review {
  id: string
  status: 'pending' | 'approved' | 'rejected'
}

function isApproved(review: Review): boolean {
  return review.status === 'approved'
}
```

## Python

```python
def is_approved(review: dict) -> bool:
    return review.get("status") == "approved"
```

## Go

```go
package main

func isApproved(status string) bool {
	return status == "approved"
}
```

## Rust

```rust
fn is_approved(status: &str) -> bool {
    status == "approved"
}
```

## Java

```java
class Review {
    static boolean isApproved(String status) {
        return "approved".equals(status);
    }
}
```

## C

```c
#include <string.h>

int is_approved(const char *status) {
    return strcmp(status, "approved") == 0;
}
```

## C++

```cpp
#include <string>

bool is_approved(const std::string &status) {
    return status == "approved";
}
```

## C#

```csharp
class Review {
    public static bool IsApproved(string status) => status == "approved";
}
```

## Ruby

```ruby
def approved?(status)
  status == "approved"
end
```

## PHP

```php
<?php
function isApproved(string $status): bool {
    return $status === "approved";
}
```

## Swift

```swift
func isApproved(_ status: String) -> Bool {
    return status == "approved"
}
```

## Kotlin

```kotlin
fun isApproved(status: String) = status == "approved"
```

## Scala

```scala
object Review {
  def isApproved(status: String): Boolean = status == "approved"
}
```

## Haskell

```haskell
isApproved :: String -> Bool
isApproved status = status == "approved"
```

## Elixir

```elixir
defmodule Review do
  def approved?(status), do: status == "approved"
end
```

## Erlang

```erlang
is_approved(Status) -> Status =:= "approved".
```

## Clojure

```clojure
(defn approved? [status]
  (= status "approved"))
```

## Lua

```lua
local function is_approved(status)
  return status == "approved"
end
```

## Perl

```perl
sub is_approved {
    my ($status) = @_;
    return $status eq "approved";
}
```

## R

```r
is_approved <- function(status) {
  status == "approved"
}
```

## SQL

```sql
SELECT id, title
FROM reviews
WHERE status = 'approved'
ORDER BY due_date;
```

## Bash

```bash
is_approved() {
  [[ "$1" == "approved" ]]
}
```

## PowerShell

```powershell
function Test-Approved {
    param([string]$Status)
    $Status -eq "approved"
}
```

## YAML

```yaml
review:
  id: r-1024
  status: approved
  reviewer: alice
```

## JSON

```json
{
  "id": "r-1024",
  "status": "approved",
  "reviewer": "alice"
}
```

## XML

```xml
<review id="r-1024">
  <status>approved</status>
  <reviewer>alice</reviewer>
</review>
```

## HTML

```html
<article class="review" data-status="approved">
  <h2>Review r-1024</h2>
  <p>Reviewer: alice</p>
</article>
```

## CSS

```css
.review[data-status="approved"] {
  border-left: 4px solid green;
}
```

## SCSS

```scss
.review {
  &[data-status="approved"] {
    border-left: 4px solid $color-success;
  }
}
```

## Markdown

```markdown
# Review r-1024

**Status:** approved
```

## Dockerfile

```dockerfile
FROM alpine:3.22
COPY review-api /usr/local/bin/review-api
CMD ["review-api"]
```

## Makefile

```makefile
build:
	go build -o review-api ./cmd/review-api

test:
	go test ./...
```

## INI

```ini
[review]
status = approved
reviewer = alice
```

## TOML

```toml
[review]
id = "r-1024"
status = "approved"
reviewer = "alice"
```

## Diff

```diff
- status: pending
+ status: approved
```

## GraphQL

```graphql
query GetReview($id: ID!) {
  review(id: $id) {
    id
    status
    reviewer
  }
}
```

## Protocol Buffers

```protobuf
message Review {
  string id = 1;
  string status = 2;
  string reviewer = 3;
}
```

## Nginx

```nginx
location /reviews/ {
  proxy_pass http://review-api:8080;
}
```

## Apache

```apache
<Location /reviews>
  ProxyPass http://review-api:8080
</Location>
```
