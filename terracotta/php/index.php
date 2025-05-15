<?php
$posts = array_reverse(glob('posts/*.txt'));
?>
<!DOCTYPE html>
<html>
<head>
  <meta charset="UTF-8">
  <title>Terracotta</title>
  <link rel="stylesheet" href="style.css">
</head>
<body>
  <h1>Terracotta</h1>

  <form action="post.php" method="post">
    <textarea name="content" rows="4" cols="50" placeholder="What's on your mind?" required></textarea><br>
    <button type="submit">Post</button>
  </form>

  <hr>

  <?php foreach ($posts as $file): ?>
    <?php $lines = file($file); ?>
    <div class="post">
      <p><?= htmlspecialchars(implode("", array_slice($lines, 1))) ?></p>
      <small><?= htmlspecialchars(trim($lines[0])) ?></small>
    </div>
    <hr>
  <?php endforeach; ?>
</body>
</html>
