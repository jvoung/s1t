load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test", "nogo")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    prefix = "github.com/jvoung/s1t",
)

go_library(
    name = "go_default_library",
    srcs = [
        "dimacs_parser.go",
        "problem_spec.go",
        "solution.go",
        "solver.go",
    ],
    importpath = "github.com/jvoung/s1t",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = [
        "dimacs_parser_test.go",
        "solver_test.go",
    ],
    data = glob([
        "test_cnf/*",
        "test_cnf_slow/*",
    ]),
    embed = [":go_default_library"],
    deps = ["@com_github_google_go_cmp//cmp:go_default_library"],
)

# For now, this might be a bit redundant with "vet" which is run for tests.
# Also, it applies to external dependencies which may not be clean.
nogo(
    name = "s1t_nogo",
    vet = True,
    visibility = ["//visibility:public"],  # must have public visibility
    deps = [
    ],
)
