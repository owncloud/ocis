/** @deprecated */
export abstract class SharePermissionBit {
  static readonly Internal: number = 0
  static readonly Read: number = 1
  static readonly Update: number = 2
  static readonly Create: number = 4
  static readonly Delete: number = 8
  static readonly Share: number = 16
}
