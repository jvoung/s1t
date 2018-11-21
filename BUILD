load("@io_bazel_rules_go//go:def.bzl", "go_library", "go_test")
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "gazelle",
    prefix = "github.com/jvoung/s1t",
)

go_library(
    name = "dimacs_parser",
    srcs = [
        "dimacs_parser.go",
        "problem_spec.go",
    ],
    importpath = "github.com/jvoung/s1t",
    visibility = ["//visibility:public"],
)

go_test(
    name = "go_default_test",
    srcs = ["dimacs_parser_test.go"],
    embed = [":go_default_library"],
    # TODO(jvoung): fix this dependency.
    deps = ["@com_github_google_go_cmp//cmp:go_default_library"],
)
