use cardano/assets.{PolicyId}
use cardano/transaction.{Transaction}

validator message_gate {
  mint(redeemer: ByteArray, _policy_id: PolicyId, _self: Transaction) {
    redeemer == "HelloMint"
  }

  else(_) {
    fail
  }
}
