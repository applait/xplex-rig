/**
 * Config sample
 */

module.exports = {

  /**
   * Server specific config
   */
  server: {

    /**
     * Host and port to run xplex-rig's server on
     */
    host: '0.0.0.0',
    port: '8081',

    /**
     * PostgreSQL DB URL
     */
    postgres_url: 'postgres://user:pass@localhost/name',

    /**
     * etcd v2 cluster information
     */
    etcd_v2: null

  },

  /**
   * Agent specific config
   */
  agent: {

    /**
     * Host and port to run xplex-rig's agent on
     */
    port: '9000',

    /**
     * URL to rig-server instance's API
     */
    rig_server_url: 'https://127.0.0.1:8081',

    /**
     * etcd v2 cluster information
     */
    etcd_v2: null,

    /**
     * Public IP address and port (optional) on which media server is connected.
     * Auto discovered if not set.
     *
     * e.g., 122.133.144.155
     */
    public_address: null,

    /**
     * Private IP address and port(optional) on which agent exposes its API.
     * Auto discovered if not set.
     *
     * e.g:, 10.11.12.13:9000
     */
    private_address: null
  }

}
