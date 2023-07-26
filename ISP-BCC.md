# Controller Testing:
### ISP - Input Space Partitioning

The main goal of input space partitioning is to simplify the problem or analysis by breaking it down into smaller, more manageable parts. By dividing the input space into partitions, it becomes possible to focus on understanding the behavior or properties of each partition separately, rather than considering the entire input space as a whole.
In the context of software testing and validation, input space partitioning is often employed to design test cases that adequately cover the various partitions of the input space. By selecting representative test inputs from each partition, our goal is to ensure that our operator  is thoroughly tested and its behavior in different scenarios is evaluated.
### BCC - Basic Choice Coverage

Base choice coverage is a testing criterion that focuses on testing the various combinations of independent options or configurations within a system. It aims to ensure that all critical combinations of base choices are considered during testing, without the need to exhaustively test all possible combinations. We will use our partitions from ISP to create our combinations for our BCC testing.

#### Reconcile()
| Parameter   | Type  | Partition    | Value | Expected Output |
|---|---|---|---|---|
| server   | *spirev1.SpireServer       | | | |
|          | server.Namespace (string)  | | | |
|          | server.Name (string)       | | | |

#### validateYaml()
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| | | | | |

#### createServiceAccount(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| | | | | |

#### spireBundleDeployment(namespace string)
| Parameter   | Type  | Partition    | Value | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| UD namespace | true |
|          |        |                    | namespace != UD namespace|  "anythingelse" | false |
|          |       | len(namespace) =< 0 | namespace == "" | "" | false |
| bundle | &corev1.ConfigMap          | typeOf(bundle) != corev1 | &rbacv1.Role  | false |
|        |                            | typeOf(bundle) == corev1 | &corev1.ConfigMap | true |
|        | bundle.Name (string)       | name == UD name| UD name | true |
|        |                            | name != UD name| "anythingelse" | false |
|        | bundle.Namespace (string)  | namespace == UD namespace| UD namespace | true |
|        |                            | namespace != UD namespace|  "anythingelse" | false |
|        | bundle.Kind (string)       | Kind != "ConfigMap" | "NotRightValue" | fase |
|        |                            | Kind == "ConfigMap" | "ConfigMap" | true |
|        | bundle.APIVersion (string) | APIVersion != "v1" | "NotRightValue" | fase |
|        |                            | APIVersion == "v1" | "v1" | true |

#### spireRoleDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| serverRole |  &rbacv1.Role | | | |
|        | Rules       | | | |
|        | serverRole.Namespace (string)  | | | |
|        | serverRole.Name (string)       | | | |
|        | serverRole.Kind (string)       | | | |
|        | serverRole.APIVersion (string) | | | |
| Rules  | rbacv1.PolicyRule     | | | |
|        | Verbs ([]string)  | | | |
|        | Resources ([]string)      | | | |
|        | APIGroups ([]string)       | | | |

#### spireRoleBindingDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| serverRole |  &rbacv1.RoleBinding | | | |
|        | RoleRef       | | | |
|        | Subjects       | | | |
|        | serverRole.Namespace (string)  | | | |
|        | serverRole.Name (string)       | | | |
|        | serverRole.Kind (string)       | | | |
|        | serverRole.APIVersion (string) | | | |
| RoleRef  | rbacv1.RoleRef     | | | |
|        | Kind (string)  | | | |
|        | Name (string)      | | | |
|        | APIGroups (string)       | | | |
| Subject  | rbacv1.Subject     | | | |
|        | Kind (string)  | | | |
|        | Name (string)      | | | |
|        | Namespace (string)      | | | |

#### spireClusterRoleDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| clusterRole |  &rbacv1.ClusterRole | | | |
|        | Rules       | | | |
|        | clusterRole.Name (string)       | | | |
|        | clusterRole.Kind (string)       | | | |
|        | clusterRole.APIVersion (string) | | | |
| Rules  | rbacv1.PolicyRule     | | | |
|        | Verbs ([]string)  | | | |
|        | Resources ([]string)      | | | |
|        | APIGroups ([]string)       | | | |


#### spireClusterRoleBindingDeployment(namespace string)
| Parameter   | Type  | Partition   | Value  | Expected Output |
|---|---|---|---|---|
| namespace| string | len(namespace) > 0 | namespace == UD namespace| true |
|           |       |                    | namespace != UD namespace| false |
|           |       | len(namespace) =< 0 | "" | false |
| serverRole |  &rbacv1.RoleBinding | | | |
|        | RoleRef       | | | |
|        | Subjects       | | | |
|        | serverRole.Name (string)       | | | |
|        | serverRole.Kind (string)       | | | |
|        | serverRole.APIVersion (string) | | | |
| RoleRef  | rbacv1.RoleRef     | | | |
|        | Kind (string)  | | | |
|        | Name (string)      | | | |
|        | APIGroups (string)       | | | |
| Subject  | rbacv1.Subject     | | | |
|        | Kind (string)  | | | |
|        | Name (string)      | | | |
|        | Namespace (string)      | | | |

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