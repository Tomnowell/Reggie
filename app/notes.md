echo -n "apple123" | ./your_program.sh -E "\d"
echo -n "alpha_num3ric" | ./your_program.sh -E "\w"
echo -n "apple" | ./your_program.sh -E "[abc]"
echo -n "apple" | ./your_program.sh -E "[^abc]"

The execution is almost correct, I think the problem lies in how the tokens are being created....