﻿// Copyright 2016-2021, Pulumi Corporation

namespace Pulumi.Automation
{
    /// <summary>
    /// Options controlling the behavior of an <see cref="WorkspaceStack.DestroyAsync(DestroyOptions, System.Threading.CancellationToken)"/> operation.
    /// </summary>
    public sealed class DestroyOptions : UpdateOptions
    {
        public bool? TargetDependents { get; set; }

        public bool? Remove { get; set; }
        /// <summary>
        /// Show config secrets when they appear.
        /// </summary>
        public bool? ShowSecrets { get; set; }
    }
}
