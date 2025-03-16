using System;

namespace DeepGate.Models.Responses
{
    public class NodeSortedHostResponse
    {
        /// <summary>
        /// The IP address of the host server
        /// </summary>
        public string HostIp { get; set; }

        /// <summary>
        /// The name/identifier of the host server
        /// </summary>
        public string HostName { get; set; }
    }
}