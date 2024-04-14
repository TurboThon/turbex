export type File = {
  id: string
  length: number
  chunkSize: number
  uploadDate: string
  senderUserName: string
  filename: string
  encryptionKey: string
  ephemeralPubKey: string
  expirationDate: string
  canWrite: boolean
}
