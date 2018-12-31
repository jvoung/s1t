load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
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
