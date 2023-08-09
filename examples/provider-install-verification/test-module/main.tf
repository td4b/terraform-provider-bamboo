
# resource "null_resource" "ad_user_test" {
#   # Changes to any instance of the cluster requires re-provisioning
#   triggers = {
#     firstname = var.firstname
#     lastname = var.lastname
#     email = var.email
#     status = var.status
#   }
# }

resource "null_resource" "ad_user_test" {
  # Changes to any instance of the cluster requires re-provisioning
  triggers = {
    values = var.values
  }
}
