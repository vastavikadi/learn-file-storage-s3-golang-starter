# Single Machine Architecture - Limitations and Trade-offs

> This lesson introduces the traditional single-machine web application architecture and explains its limitations (scaling, availability, durability, cost, maintenance) that motivate the need for serverless solutions like AWS S3.

## Prerequisites
- Basic understanding of web application architecture
- Familiarity with cloud computing concepts (VMs, instances)
- Knowledge of HTTP servers, databases, and file systems

## Key Concepts

### Single Machine Architecture

A simple web application architecture where all components run on one physical or virtual machine in the cloud. This includes the HTTP server, database, and file system all co-located on the same instance.

**Examples:**
- A typical EC2 instance running nginx/Apache, PostgreSQL/MySQL, and storing uploaded files on local disk
- A VPS (DigitalOcean Droplet, Linode, Vultr) running the entire application stack
- A developer's laptop running everything locally during development

> ⚠️ **Gotchas:**
> - While simple to set up, this architecture creates a single point of failure
> - All resources (CPU, RAM, disk, network) are shared and contend with each other
> - Vertical scaling has hard physical limits

### Components of a Single Machine

The three core components typically running on a single machine architecture:

**Examples:**
- HTTP Server: Handles incoming requests (nginx, Apache, Go http.Server, Node.js Express)
- Database: Stores structured data (PostgreSQL, MySQL, MongoDB, SQLite)
- File System: Stores larger unstructured files (user uploads, generated reports, static assets)

> ⚠️ **Gotchas:**
> - Database and file storage compete for disk I/O
> - HTTP server worker processes compete with database for CPU/RAM
> - No isolation - a memory leak in one component affects all others

### Scaling Limitations (Vertical Scaling Only)

Single machine architectures can only scale vertically (adding more resources to the same machine). There's a hard ceiling on how powerful a single machine can become.

**Examples:**
- AWS EC2 instance types max out at certain CPU/RAM (e.g., u-24tb1.metal with 24 TB RAM, 448 vCPUs)
- Disk space limited by attached EBS volumes (max 64 TB per volume, but practical limits lower)
- Network bandwidth capped by instance type

> ⚠️ **Gotchas:**
> - Vertical scaling requires downtime (stop instance, change type, start)
> - Cost grows exponentially at higher tiers
> - Diminishing returns - doubling CPU doesn't double throughput due to contention

### Availability Concerns

Single machine = single point of failure. If the machine goes down (hardware failure, kernel panic, OOM killer, network partition), the entire application is unavailable.

**Examples:**
- AWS EC2 instance retirement due to hardware degradation
- Kernel panic from a buggy driver or memory corruption
- OOM killer terminating the database process
- Network partition isolating the instance from clients

> ⚠️ **Gotchas:**
> - Load balancers + multiple servers mitigate this, but that's no longer a 'single machine' architecture
> - Even with auto-recovery, there's minutes of downtime during instance replacement
> - Stateful components (database, file system) make horizontal scaling complex

### Durability Risks

Data stored on a single machine's local disk is vulnerable to loss from hardware failure, human error, or software bugs. Backups are often neglected or inadequate.

**Examples:**
- Disk failure (SSD wear-out, controller failure) - RAID helps but isn't foolproof
- Accidental `rm -rf /var/lib/mysql` or `rm -rf /app/uploads`
- Buggy deployment script that wipes the data directory
- Ransomware or malicious insider deleting data
- EBS volume corruption (rare but possible)

> ⚠️ **Gotchas:**
> - Snapshots/backups have RPO (Recovery Point Objective) - data since last backup is lost
> - Restoring from backup takes time (RTO - Recovery Time Objective)
> - Most developers don't test restore procedures regularly
> - Local disk on EC2 is ephemeral (instance store) - data lost on stop/terminate

### Cost Inefficiency

Running a server 24/7 means paying for idle capacity. Most applications have variable load (day/night cycles, weekends, seasonal spikes) but single machine architecture requires provisioning for peak.

**Examples:**
- t3.medium at ~$30/month running 24/7 = $360/year even if only used 8 hours/day
- Over-provisioned for Black Friday traffic, idle rest of year
- Database connection pooling limits mean you need larger instance than CPU/RAM alone would suggest

> ⚠️ **Gotchas:**
> - Reserved Instances / Savings Plans reduce cost but lock you in
> - Spot instances can reduce cost 90% but can be terminated with 2-minute notice
> - Right-sizing is a continuous optimization problem

### Operational Burden (Maintenance)

Running your own server means you're responsible for all operational tasks - the 'ops' in DevOps. This distracts from core product development.

**Examples:**
- OS patching and security updates (kernel, libc, OpenSSL, etc.)
- Database version upgrades (major version migrations are risky)
- Log rotation, monitoring, alerting setup
- Backup scheduling, verification, and restore testing
- SSL certificate renewal (Let's Encrypt every 90 days)
- Disk space monitoring and cleanup
- Security hardening (fail2ban, ufw, SSH key rotation)

> ⚠️ **Gotchas:**
> - Toil work that doesn't directly add customer value
> - Requires specialized knowledge (Linux sysadmin, DBA skills)
> - Bus factor - if the ops person leaves, knowledge is lost
> - Compliance requirements (SOC2, HIPAA) add significant overhead

## Mental Models

### 💡 The Apartment Building vs. Single Family Home

Single machine = single family home. You own the whole thing, but if the roof leaks, you're homeless. Serverless/S3 = apartment building with professional maintenance. You rent one unit; the landlord handles roof, plumbing, electricity. You share infrastructure but get reliability.

**Analogy:** Single machine: You're the homeowner, plumber, electrician, and security guard. Serverless: You're a tenant paying rent for guaranteed utilities.

### 💡 The Egg Basket Analogy

Putting all your eggs (HTTP, DB, files) in one basket (single machine). Drop the basket, lose all eggs. Distributed systems = multiple baskets. S3 = a specialized, industrial-grade egg vault with 99.999999999% durability guarantee.

**Analogy:** Single machine = one basket carried by one person. S3 = Fort Knox for eggs.

### 💡 The Restaurant Kitchen

Single machine = one chef doing everything: taking orders, cooking, washing dishes, managing inventory, cleaning. Works for a food truck. Fails at scale. S3 = dedicated cold storage warehouse. Chef focuses on cooking; warehouse handles storage at massive scale.

**Analogy:** Don't make your web server also be your file warehouse. Specialization wins at scale.

## Code Examples

### Typical Single Machine Docker Compose (Development)

Even in development, a single-machine architecture often uses multiple containers on one host. This simulates production but shares the same underlying resources.

```yaml
version: '3.8'
```

### Single Machine Systemd Service File

A typical systemd unit file for running a Go web server on a single machine. The same machine would also run PostgreSQL and nginx via separate unit files.

```ini
[Unit]
```

### Backup Script (The Thing You Probably Don't Have)

A basic backup script that many teams intend to write but never get around to. This illustrates the durability gap in single-machine architectures.

```bash
#!/bin/bash
```

### Vertical Scaling Downtime Example

Vertical scaling requires stopping the instance, changing the instance type, and starting it again - causing downtime.

```bash
# AWS CLI commands to vertically scale an EC2 instance
```

### Disk Space Monitoring Alert

Operational toil: you need to monitor disk space yourself on a single machine. In S3, this is handled automatically.

```bash
# Cron job to alert on disk usage
```

## Common Mistakes

- ❌ Assuming single machine is 'good enough' for production without considering bus factor
- ❌ Not having tested backup/restore procedures (backup != restore capability)
- ❌ Thinking vertical scaling is infinite - hitting instance type limits unexpectedly
- ❌ Running database and file storage on same disk - I/O contention kills performance
- ❌ Using instance store (ephemeral) storage for persistent data on EC2
- ❌ Neglecting OS/security patches until a CVE forces emergency maintenance
- ❌ Not monitoring disk space until the server crashes with 'no space left on device'
- ❌ Believing 'it works on my machine' translates to production reliability

## Key Takeaways

- ✅ Single machine architecture is valid for simple apps, prototypes, and low-traffic production workloads
- ✅ Five fundamental trade-offs: scaling ceiling, availability (SPOF), durability risk, 24/7 cost, operational burden
- ✅ Vertical scaling has hard limits and requires downtime; horizontal scaling requires architectural changes
- ✅ Durability on a single machine requires disciplined backup strategy that most teams don't implement
- ✅ Operational toil (patching, monitoring, backups, upgrades) distracts from product development
- ✅ S3 addresses the file storage component of this architecture with serverless, durable, scalable, pay-per-use storage
- ✅ Understanding these limitations is why S3's 'serverless' model is compelling - it eliminates the file system component from your server

## Practice Questions

**Q1:** What are the three core components typically running on a single machine web application architecture?

<details><summary>Answer</summary>

HTTP server, database, and file system (for larger files).

</details>

**Q2:** Why does vertical scaling on a single machine require downtime?

<details><summary>Answer</summary>

Changing instance types (CPU/RAM) on cloud providers like AWS requires stopping the instance, modifying the instance type, and starting it again. The OS and applications must restart.

</details>

**Q3:** What is the durability risk of storing files on a single machine's local disk?

<details><summary>Answer</summary>

Data can be lost due to hardware failure (disk death), human error (accidental deletion), software bugs, or ransomware. Backups mitigate but have RPO/RTO gaps and are often untested.

</details>

**Q4:** How does the cost model of a single machine differ from serverless storage like S3?

<details><summary>Answer</summary>

Single machine: pay 24/7 for provisioned capacity regardless of usage. S3: pay per GB stored per month + per request, with no provisioning. S3 scales to zero cost when unused.

</details>

**Q5:** What operational tasks does a single machine architecture require that S3 eliminates for file storage?

<details><summary>Answer</summary>

Disk space monitoring, backup scheduling/verification, RAID management, filesystem corruption checks, OS patching for storage drivers, capacity planning for disk growth, log rotation for storage services.

</details>

## Concept Relationships

- **Single Machine Architecture** → **AWS S3**: S3 replaces the file system component of single machine architecture with a serverless, durable, scalable alternative
- **Single Machine Architecture** → **Load Balancers + Multiple Servers**: Horizontal scaling solution that addresses availability but adds complexity and doesn't solve file storage sharing
- **Single Machine Architecture** → **Database (RDS, DynamoDB)**: Managed database services address the database component similarly to how S3 addresses file storage
- **Single Machine Architecture** → **Container Orchestration (ECS, Kubernetes)**: Modern way to run multiple machines with shared infrastructure, but still requires managing the control plane

---
*Generated by Boot.dev Notes*