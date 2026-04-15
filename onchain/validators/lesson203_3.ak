use cardano/assets.{PolicyId}
use cardano/transaction.{Transaction}

validator lesson203_3 {
  mint(_redeemer: Data, _policy_id: PolicyId, _self: Transaction) {
    True
  }

  else(_) {
    fail
  }
}