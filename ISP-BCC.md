# Controller Testing:
### ISP - Input Space Partitioning

The main goal of input space partitioning is to simplify the problem or analysis by breaking it down into smaller, more manageable parts. By dividing the input space into partitions, it becomes possible to focus on understanding the behavior or properties of each partition separately, rather than considering the entire input space as a whole.
In the context of software testing and validation, input space partitioning is often employed to design test cases that adequately cover the various partitions of the input space. By selecting representative test inputs from each partition, our goal is to ensure that our operator  is thoroughly tested and its behavior in different scenarios is evaluated.
### BCC - Basic Choice Coverage

Base choice coverage is a testing criterion that focuses on testing the various combinations of independent options or configurations within a system. It aims to ensure that all critical combinations of base choices are considered during testing, without the need to exhaustively test all possible combinations. We will use our partitions from ISP to create our combinations for our BCC testing.

### Reconcile(ctx context.Context, req ctrl.Request)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| server | *spirev1.SpireServer       | | | |
|        | server.Namespace (string)  | | | |
|        | server.Name (string)       | | | |

#### validateYaml()
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| | | | | |

#### createServiceAccount(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| | | | | |

#### spireBundleDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| bundle | &corev1.ConfigMap | | | |
|        | bundle.Name (string)       | | | |
|        | bundle.Namespace (string)  | | | |
|        | bundle.Name (string)       | | | |
|        | bundle.Kind (string)       | | | |
|        | bundle.APIVersion (string) | | | |

#### spireRoleDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| serverRole |  &rbacv1.Role | | | |
|        | Rules       | | | |
|        | serverRole.Namespace (string)  | len(namespace) > 0 | namespace == UD namespace| true |
|        |                                |                    | namespace != UD namespace| false |
|        |                                | len(namespace) =< 0 | "" | false |
|        | serverRole.Name (string)       | name == "spire-server-configmap-role"| "spire-server-configmap-role"| true |
|        |                                | name != "spire-server-configmap-role"| "anythingElse"| false |
|        | serverRole.Kind (string)       | Kind == "Role" | "Role" |true |
|        |                                | Kind != "Role" | "anythingElse" |false |
|        | serverRole.APIVersion (string) | APIVersion == "rbac.authorization.k8s.io/v1"| "rbac.authorization.k8s.io/v1"| true |
|        |                                | APIVersion != "rbac.authorization.k8s.io/v1"| "anythingElse"| false |
| Rules  | rbacv1.PolicyRule     | | | |
|        | Verbs ([]string) | len(Verbs) >= 3 | []string{"patch", "get", "list"}| true |
|        |                  |                 | []string{"addams", "get", "list"}| false |
|        |                  | len(Verbs) < 3  | []string{"get", "list"} | false|
|        | Resources ([]string) | len(Resources) == 1 | []string{"configmaps"} | true |
|        |                      |                     | []string{"whatnot"} | false |
|        |                      | len(Resources) != 1 | []string{}| false |
|        | APIGroups ([]string) | len(APIGroups) == 1 | []string{""}| true |
|        |                      |                     | []string{"blah"}| false |
|        |                      | len(APIGroups) != 1 | []string{}| false |

#### spireRoleBindingDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| serverRole |  &rbacv1.RoleBinding | | | |
|        | RoleRef       | | | |
|        | Subjects       | | | |
|        | serverRole.Namespace (string)  | len(namespace) > 0 | namespace == UD namespace| true |
|        |                                |                    | namespace != UD namespace| false |
|        |                                | len(namespace) =< 0 | "" | false |
|        | serverRole.Name (string)       | name == "spire-server-configmap-role-binding"| "spire-server-configmap-role-binding"| true |
|        |                                | name != "spire-server-configmap-role-binding"| "anythingElse"| false |
|        | serverRole.Kind (string)       | Kind == "RoleBinding" | "RoleBinding" |true |
|        |                                | Kind != "RoleBinding" | "anythingElse" |false |
|        | serverRole.APIVersion (string) | APIVersion == "rbac.authorization.k8s.io/v1"| "rbac.authorization.k8s.io/v1"| true |
|        |                                | APIVersion != "rbac.authorization.k8s.io/v1"| "anythingElse"| false |
| RoleRef  | rbacv1.RoleRef     | | | |
|        | Name (string) | name == "spire-server-configmap-role"| "spire-server-configmap-role"| true |
|        |               | name != "spire-server-configmap-role"| "anythingElse"| false |
|        | Kind (string) | Kind == "Role" | "Role" |true |
|        |               | Kind != "Role" | "anythingElse" |false |
|        | APIVersion (string) | APIVersion == "rbac.authorization.k8s.io/v1"| "rbac.authorization.k8s.io/v1"| true |
|        |                     | APIVersion != "rbac.authorization.k8s.io/v1"| "anythingElse"| false |
| Subject  | rbacv1.Subject     | | | |
|        | Kind (string)     | Kind == "ServiceAccount"| "ServiceAccount" | true |
|        |                   | Kind != "ServiceAccount"|"anythingElse" | false |
|        | Name (string)     |len(name) > 0 | name == "spire-server" | true |
|        |                   |                    | namespace != "spire-server"| false |
|        |                   | len(name) =< 0 | "" | false |
|        | Namespace (string)| len(namespace) > 0 | namespace == UD namespace| true |
|        |                   |                    | namespace != UD namespace| false |
|        |                   | len(namespace) =< 0 | "" | false |

#### spireClusterRoleDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| clusterRole |  &rbacv1.ClusterRole | | | |
|        | Rules       | | | |
|        | clusterRole.Name (string)       | name == "spire-server-trust-role"| "spire-server-trust-role"| true |
|        |                                 | name != "spire-server-trust-role"| "anythingElse"| false |
|        | clusterRole.Kind (string)       | Kind == "ClusterRole" | "ClusterRole" |true |
|        |                                 | Kind != "ClusterRole" | "anythingElse" |false |
|        | clusterRole.APIVersion (string) | APIVersion == "rbac.authorization.k8s.io/v1"| "rbac.authorization.k8s.io/v1"| true |
|        |                                 | APIVersion != "rbac.authorization.k8s.io/v1"| "anythingElse"| false |
| Rules  | rbacv1.PolicyRule     | | | |
|        | Verbs ([]string) | len(Verbs) == 1 | []string{"create"} | true |
|        |                  |                 | []string{"addams", "get", "list"}| false |
|        |                  | len(Verbs) != 1 | []string{"get", "list"} | false|
|        | Resources ([]string) | len(Resources) == 1 | []string{"tokenreviews"} | true |
|        |                      |                     | []string{"whatnot"} | false |
|        |                      | len(Resources) != 1 | []string{}| false |
|        | APIGroups ([]string) | len(APIGroups) == 1 | []string{"authentication.k8s.io"} | true |
|        |                      |                     | []string{"blah"}| false |
|        |                      | len(APIGroups) != 1 | []string{}| false |



#### spireClusterRoleBindingDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| clusterRoleBinding |  &rbacv1.RoleBinding | | | |
|        | RoleRef       | | | |
|        | Subjects       | | | |
|        | clusterRoleBinding.Name (string)       | name == "spire-server-trust-role-binding"| "spire-server-trust-role-binding"| true |
|        |                                        | name != "spire-server-trust-role-binding"| "anythingElse"| false |
|        | clusterRoleBinding.Kind (string)       | Kind == "ClusterRoleBinding" | "ClusterRoleBinding" |true |
|        |                                        | Kind != "ClusterRoleBinding" | "anythingElse" |false |
|        | clusterRoleBinding.APIVersion (string) | APIVersion == "rbac.authorization.k8s.io/v1"| "rbac.authorization.k8s.io/v1"| true |
|        |                                        | APIVersion != "rbac.authorization.k8s.io/v1"| "anythingElse"| false |
| RoleRef  | rbacv1.RoleRef     | | | |
|        | Name (string)     | name == "spire-server-trust-role"| "spire-server-trust-role"| true |
|        |                   | name != "spire-server-trust-role"| "anythingElse"| false |
|        | Kind (string)     | Kind == "ClusterRole" | "ClusterRole" |true |
|        |                   | Kind != "ClusterRole" | "anythingElse" |false |
|        | APIGroup (string) | APIVersion == "rbac.authorization.k8s.io/v1"| "rbac.authorization.k8s.io/v1"| true |
|        |                   | APIVersion != "rbac.authorization.k8s.io/v1"| "anythingElse"| false |
| Subject  | rbacv1.Subject     | | | |
|        | Kind (string)     | Kind == "ServiceAccount"| "ServiceAccount" | true |
|        |                   | Kind != "ServiceAccount"|"anythingElse" | false |
|        | Name (string)     |len(name) > 0 | name == "spire-server" | true |
|        |                   |                    | namespace != "spire-server"| false |
|        |                   | len(name) =< 0 | "" | false |
|        | Namespace (string)| len(namespace) > 0 | namespace == UD namespace| true |
|        |                   |                    | namespace != UD namespace| false |
|        |                   | len(namespace) =< 0 | "" | false |

#### spireConfigMapDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| | | | | |

#### spireStatefulSetDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| | | | | |

#### spireServiceDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| | | | | |