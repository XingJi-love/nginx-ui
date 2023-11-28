import type { ModelBase } from '@/api/curd'
import Curd from '@/api/curd'
import type { DNSProvider } from '@/api/auto_cert'

export interface DnsCredential extends ModelBase {
  name: string
  config?: DNSProvider
  provider: string
}

const dns_credential: Curd<DnsCredential> = new Curd('/dns_credential')

export default dns_credential
