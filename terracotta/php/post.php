<?php
if ($_SERVER['REQUEST_METHOD'] === 'POST' && !empty($_POST['content'])) {
    $content = trim($_POST['content']);
    $timestamp = date('Y-m-d H:i:s');
    $filename = 'posts/' . time() . '.txt';

    file_put_contents($filename, $timestamp . PHP_EOL . $content);
}

header('Location: index.php');
exit;
