# Inform: Worker Updates

### Previous Worker Architecture

##### Performance

![](C:\Users\Maldonado\AppData\Roaming\marktext\images\2026-05-07-13-14-37-image.png)

- **Average Time:** 5837 ms

- **Test:** Single Submission Execution

![](C:\Users\Maldonado\AppData\Roaming\marktext\images\2026-05-07-13-11-33-image.png)

- **Average Time:** 28368 ms
- **Test:** Concurrent Submissions (n=10)

### New Worker Architecture

##### Performance

**Changes:** Implemented handling for concurrent requests

![](C:\Users\Maldonado\AppData\Roaming\marktext\images\2026-05-07-08-21-04-image.png)

- **Average Time:** 9128 ms

- **Test:** Concurrent Submissions (n=10)



**Changes:** Implemented containers pooling

![](C:\Users\Maldonado\AppData\Roaming\marktext\images\2026-05-07-11-35-26-image.png)

- **Average Time:** 5939 ms

- **Test:** Single Submission Execution

![](C:\Users\Maldonado\AppData\Roaming\marktext\images\2026-05-07-11-37-32-image.png)

- **Average Time:** 6702 ms

- **Test:** Concurrent Submissions (n=10)

![](C:\Users\Maldonado\AppData\Roaming\marktext\images\2026-05-07-11-40-16-image.png)

- **Average Time:** 9284 ms

- **Test:** Concurrent Submissions (n=20)
