## Traditional File Storage
### "File storage" is what you're already familiar with:

- Files are stored in a hierarchy of directories
- A file's system-level metadata (like timestamp and permissions) is managed by the file system, not the file itself
- File storage is great for single-machine-use (like your laptop), but it doesn't distribute well across many servers.
- It's optimized for low-latency access to a small number of files on a single machine.

## How Object Storage Differs
### Object storage is designed to be more scalable, available, and durable than file storage because it can be easily distributed across many machines:

- Objects are stored in a flat namespace (no directories)
- An object's metadata is stored with the object itself


## File System Illusion
- Object storage is flat, but it can feel like a file system. You can create the illusion of directories by using prefixes in your object keys.

- Directories are really great for organizing stuff. Storing everything in one giant bucket makes a big hard-to-manage mess. So, S3 makes your objects feel like they're in directories, even though they're not.

### It's Just Prefixes
- Keys inside of a bucket are just strings.
- If you upload an object to S3 with the key users/john/profile.jpg, we can kind of pretend that the object is in a directory called users and a subdirectory called john.

```
users/dan/profile.jpg
users/dan/friends.jpg
users/lane/profile.jpg
users/lane/friends.jpg
people/matt/profile.jpg
```

- Although directories are an illusion in S3, they're still useful due to the prefix filtering capabilities of the S3 API.
- There are a lot of common strategies for organizing objects in S3.

- We always want to group objects in a way that makes sense for our case, because often we'll want to operate on a group of objects at once.

```
For example, pretend you do the naive thing and upload all your images to the root of your bucket. What happens if...

you want to delete all the images for a specific user?
a feature changed and you need to resize all the images it uses?
you want to change the permissions of all the images associated with a specific organization?

If you don't have any prefixes (directories) to group objects, you might find yourself iterating over every object in the bucket to find the ones you care about. That's slow and expensive.
```