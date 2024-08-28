package examples

type LanguageCIExample struct {
	Language string `json:"language"`
	Example  string `json:"example"`
}

const (
	GolangExample = `
Reference: https://github.com/gotestyourself/gotestsum
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_reports.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html#go

test:
  stage: test
  image: golang:1.23
  script:
    - go install gotest.tools/gotestsum@latest
    - gotestsum --junitfile report.xml --format testname
  artifacts:
    when: always
    paths:
      - "report.xml"
    reports:
      junit: report.xml
`
	// lint:
	//   stage: lint
	//   image: golangci/golangci-lint:v1.57-alpine
	//   script:
	//     # Use default .golangci.yml file from the image if one is not present in the project root.
	//     - "[ -e .golangci.yml ] || cp /golangci/.golangci.yml ."
	//     - golangci-lint run --issues-exit-code 0 --print-issued-lines=false --out-format code-climate:gl-code-quality-report.json,line-number
	//   artifacts:
	//     reports:
	//       codequality: gl-code-quality-report.json
	//     paths:
	//       - gl-code-quality-report.json
	// 		`
	JavaExample = `
Reference: https://howtodoinjava.com/junit5/xml-reports/
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_reports.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html#maven

test:
  stage: test
  image: maven:3.9.8
  script:
    - mvn test
  artifacts:
  	when: always
  	paths:
	  	- "target/surefire-reports"
  	reports:
	  	junit: target/surefire-reports/TEST-*.xml
`
	KotlinExample = `
Reference: https://kotest.io/docs/extensions/junit_xml.html
Reference: https://stackoverflow.com/questions/62527103/generate-xml-report-for-android-test-task
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_reports.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html

test:
  stage: test
  image: gradle:8.10.0
  script:
    - gradle test
  artifacts:
	when: always
	paths:
		- "build/test-results/test"
	reports:
		junit: build/test-results/test/*.xml
`
	JavaScriptExample = `
Reference: https://vitest.dev/guide/reporters
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_reports.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html

test:
  stage: test
  image: node:20.9.0
  script:
	- npx vitest --reporter=junit --outputFile=report.xml --run
  artifacts:
	when: always
	paths:
		- "report.xml"
	reports:
		junit: report.xml
`

	PythonExample = `
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_reports.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html#python-example

test:
  stage: test
  script:
    - pytest --junitxml=report.xml
  artifacts:
    when: always
	paths:
		- "report.xml"
    reports:
      junit: report.xml
`

	ElixirExample = `
Reference: https://hexdocs.pm/junit_formatter/readme.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_reports.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html

test:
  stage: test
  image: elixir:1.17.2
  script:
    - mix test
  artifacts:
	when: always
	paths:
	  - "_build/test/lib/<project-name>"
	reports:
	  junit: "_build/test/lib/**/report_file.xml"
`

	DefaultExample = `
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_reports.html
Reference: https://docs.gitlab.com/ee/ci/testing/unit_test_report_examples.html

Your language has no explicit example. Visit the Reference links for more information. And your CI configuration might look like this in the end:

test:
  stage: test
  image: <docker-image>
  script:
	- <test-command that generates JUnit XML report>
  artifacts:
	when: always
	paths:
		- "<path-to-junit-xml-report file or folder>"
	reports:
		junit: "<path-to-junit-xml-report file or folder(*.xml)>"
`
)

var languageCIExamples = map[string]string{
	"Go":         GolangExample,
	"Java":       JavaExample,
	"Kotlin":     KotlinExample,
	"JavaScript": JavaScriptExample,
	"TypeScript": JavaScriptExample,
	"Python":     PythonExample,
	"Elixir":     ElixirExample,
}

func GetLanguageCIExample(language string) LanguageCIExample {
	example, ok := languageCIExamples[language]
	if !ok {
		example = DefaultExample
	}
	return LanguageCIExample{
		Language: language,
		Example:  example,
	}
}
