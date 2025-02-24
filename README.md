# cast
A CLI tool to enable easy HTTP requests to be completed via the CLI or file-based workflows.

**Using It:**
```bash
cast post https://example.com/api/v1/userUpdate -B newName=Greg -H X-Api-Key:12320f-r434rf-g3tg45
```

HTTP and OpenAPI are great standards. **This tool revolves around *just* them.**

The focus of the CLI input and file-based scripting is to keep it simple and obvious, which requires being a little explicit. Things like:

- Methods are required. We don't choose a default for you.
- The protocol is required. We don't assume what you're using.
- Past that, the order doesn't matter, but it is broken up by flags to ensure it's obvious you're putting something in a header, vs. a body, vs. calling a different argument.

## install

### pkg repos

#### macOS

- homebrew: `brew install cast`

#### linux

- pacman (arch via aur): ``

### build it

**With gorealeaser:**

**With go:**

## example(s) / usage

### cli usage

```bash
cast [ method ] [ url ] <args>
```

- anything after the url isn't required.
- the method *then* url order cannot change.
- output is *just* the full first response by default.

**GET Request with Custom Headers and Syntax Highlighting:**

```bash
cast get https://google.com -H test:test1 -H test2:test3 --highlight
```

![image](./documentation/example_screenshots/readme_example_1.png)

### file usage

Requests are separated mainly by `<% operations %>` and `methods`. They also can be separated by asserts, as it's something the tool matches on, but **asserts are not required**. So after a requests body it'll keep looking and won't start a new request until it finds a `{{ method }} {{ endpoint }}`.

Comments can be included before methods, or after asserts. If you don't want to assert but need to leave comments directly after a body, just do: `<% null %>` and leave the comments after it. A clear endpoint for the body is important.

Asserts are separated from each other by newlines, and from the body by two newlines. If, for whatever reason, the exact string `ASSERT` is needed in the body after two newlines, it can be escaped: `ASS\ERT`. The backslash will be auto-removed by the program post-assertions checks for any string that exactly matches that.

A `Host` header **is required** for clarity. Even though in a normal HTTP interaction it technically wouldn't be required, for the scripts it is.

- `save` works to temporarily save a value, and doubles as an `assert`.
- `[function] [location] <modifier> [value] <modifier>` is a good way to see the post-response functionality.
- If one default is defined on the first use of one variable, it'll count as the default for the rest of its uses.
- If a variable is hardcoded, it'll will act as the default.
- Content-type will be auto-assigned if one isn't entered.

```bash
### Authentication Sequence
base_url = "https://api.example.com/v1"
uuid = "d2342-d4d23e-d232d3-df2f4f2"
env_token = env.get("API_TOKEN")

POST /auth HTTP/1.1
Host: {{ base_url }}
Content-Type: application/json
X-Request-ID: {{ uuid }}

{
  "user": "{{ user || default = "test")}}",
  "pass": "{{ env_token }}"
}

<% assert status 200 %>
<% assert header "Location" %>
<% assert header NOT "X-Rate-Limit" %>
<% save body "$.token" AS auth_token %>

### Resource Creation
GET /resources HTTP/1.1
Host: {{ base_url }}
Authorization: Bearer {{auth_token}}
```



## quick comparisons

### with httpie


### with hurl


### with xh


### with cURL


### benchmark


## issues, contributions, faqs

*Warning: Not all issues or feature requests may be accepted.*

- `.md` on how to [submit an issue]()
- `.md` on how to [contribute]()
- `.md` containing [the faq]()

Generally the program *shouldn't* load / initiate / do anything that isn't involved in **only** the thing the user commands it. And the things it *does* do, shouldn't be done until they're needed.

## roadmap

*Warning: Any version under v1.0 is going to be **very** unstable and subject to change (flags, syntax, etc.). Please report any bugs or issues.*

### cli focus (<= v0.5)

**v0.1**

- [x] HTTP methods + URL
- [x] Custom headers and request body

**v0.2**

- [x] Print request if debug flag is `true`.
- [ ] Track duration between request and response
- [x] Improved / useful logging with values being interacted with
- [ ] `-UF` or `--uploadfile` to place a file in the request body
- [x] Short response support. Add a print option that just prints the response's status `200 OK`.

**v0.3**

- [ ] `-R` flag to print request as well
- [ ] `--print` options similar to httpie: https://httpie.io/docs/cli/what-parts-of-the-http-exchange-should-be-printed
- [ ] Flag to allow following redirects and only returning the "final" response.

**v0.4**

- [ ] Fix error returning structure and handle errors correctly and in better places.
- [ ] Create tests & benchmark tests for all files.
- [x] Custom variables
- [ ] Cross-run variable storage / .env support

**v0.5**

- [ ] Optimize. Review for helpful debug logging (esp. at the beginning, during, and before return).
- [ ] Make the repo public and set up `goreleaser`.
- [ ] Setup packages with `brew` & `pacman`.

### file focus (> v0.5)

- [ ] File input
- [ ] Assertions

**v0.6**

- [ ] Config
- [ ] Run a directory alphabetically option `cast --directory /tests/auth` (alias `-DIR`)

**v0.7**

- [ ] Support proxying
- [ ] Support client certificates
- [ ] Export as cURL command(s) `cast --curl-export -F auth.cast`
- [ ] Export as OpenAPI v3.x (whatever is latest) `cast --openapi-export -DIR /tests/auth`

**v0.8**

- [ ] Reference files to include in the body via their path
- [ ] Wordlist / fuzzing support via the provided vars (`<% varName = ./tests/wordlist.txt %>` or something similar. Will resend the request enough times to iterate through the list. Can only be used in the request it is defined above.)

**v0.9**

- [ ] Create the docs site using Astro Starlight.
- [ ] Get feedback from people who would make use of it.
- [ ] Syntax highlighting for editing in Zed or VSCode.

**v1.0**

- [ ] Create issue templates in GitHub and check labels
- [ ] Ensure GitHub workflows are configured correctly and packaging is essentially automated
- [ ] Run a SAST in the pipeline. Something like Semgrep. Get it plugged into github issues
- [ ] Submit application to be an OWASP project (might require a license change?)
- [ ] Announcement posts. Reached the point of having enough features to be useful

### continuing (>= v1.0)
- [ ] optimize. find anti-patterns.
- [ ] maintain. keep on top of package and language changes.

### *potential* future plans
- [ ] grpc support

## writings about creating it

I wrote a few blog posts on my personal blog about my experience and thoughts while developing `cast`. The are stream-of-thought and show how the scope of the project and its features changed over time. They can be read [here](), and exist under the `cast` tag.

## license
