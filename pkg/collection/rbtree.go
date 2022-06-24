package collection

import (
	"fmt"
)

type dir string

const (
	leftDir  dir = "left"
	rightDir dir = "right"
	LLDir        = leftDir + leftDir
	LRDir        = leftDir + rightDir
	RRDir        = rightDir + rightDir
	RLDir        = rightDir + leftDir
	nilDir   dir = "nil"
)

// RBTree red black tree
// refer: https://github.com/chrislessard/LSM-Tree
type RBTree struct {
	count int
	root  *node
}

func NewRBTree() *RBTree {
	return &RBTree{
		count: 0,
		root:  nil,
	}
}

func (t *RBTree) Insert(key string, value any) {
	if t.root == nil {
		t.root = newNode(key, value, blackNode, nil, left(emptyLeaf), right(emptyLeaf))
		t.count++
		return
	}
	parent, nodeDir := t.findParent(key)
	if nodeDir == nilDir {
		parent.value = value
		return
	}
	n := newNode(key, value, redNode, parent, left(emptyLeaf), right(emptyLeaf))
	if nodeDir == leftDir {
		parent.left = n
	} else {
		parent.right = n
	}

	t.tryBalance(n)
	t.count++
}

func (t *RBTree) Remove(key string) {
	nodeToRemove := t.findNode(key)
	if nodeToRemove == nil {
		//  node is not in the tree
		return
	}
	if nodeToRemove.getChildrenCount() == 2 {
		// find the in-order successor and replace its value.
		// then, remove the successor
		successor := t.findInOrderSuccessor(nodeToRemove)
		nodeToRemove.key = successor.key // switch the value
		nodeToRemove = successor
	}
	// has 0 or 1 children!
	t.removeNode(nodeToRemove)
	t.count -= 1
}

func (t *RBTree) Get(key string) any {
	return t.findNode(key).value
}

func (t *RBTree) Contains(key string) bool {
	return t.findNode(key) != nil
}

func (t *RBTree) Ceil(key string) any {
	if t.root == nil {
		return nil
	}
	lastVal := t.root
	if t.root.key < key {
		lastVal = nil
	}
	t.findCeil(t.root, lastVal, key)
	return lastVal.value
}

func (t *RBTree) findCeil(n, lastVal *node, key string) {
	if n.equals(emptyLeaf) {
		return
	}
	if n.key == key {
		lastVal = n
		return
	}
	if n.key < key {
		t.findCeil(n.right, lastVal, key)
	}
	lastVal = n
	t.findCeil(n.left, lastVal, key)
}

func (t *RBTree) Floor(key string) any {
	if t.root == nil {
		return nil
	}
	lastVal := t.root
	if t.root.key > key {
		lastVal = nil
	}

	t.findFloor(t.root, lastVal, key)

	return lastVal.value
}

func (t *RBTree) findFloor(n, lastVal *node, key string) {
	if n.equals(emptyLeaf) {
		return
	}
	if n.key == key {
		lastVal = n
		return
	}
	if n.key < key {
		lastVal = n
		t.findFloor(n.right, lastVal, key)
	}
	t.findFloor(n.left, lastVal, key)
}

func (t *RBTree) removeNode(n *node) {
	leftChild := n.left
	rightChild := n.right
	notEmptyChild := rightChild
	if !leftChild.equals(emptyLeaf) {
		notEmptyChild = leftChild
	}
	if n.equals(t.root) {
		if !notEmptyChild.equals(emptyLeaf) {
			t.root = notEmptyChild
			t.root.parent = nil
			t.root.color = blackNode
		} else {
			t.root = nil
		}
	} else if n.color == redNode {
		if !n.hasChildren() {
			t.removeLeaf(n)
		} else {
			/*
				Since the node is red he cannot have a child.
				If he had a child, it'd need to be black, but that would mean that
				the black height would be bigger on the one side and that would make our tree invalid
			*/
			panic("Unexpected behavior")
		}
	} else {
		if rightChild.hasChildren() || leftChild.hasChildren() {
			panic("The red child of a black node with 0 or 1 children" +
				" cannot have children, otherwise the black height of the tree becomes invalid! ")
		}
		if notEmptyChild.color == redNode {
			n.key = notEmptyChild.key
			n.left = notEmptyChild.left
			n.right = notEmptyChild.right
		} else {
			t.removeBlackNode(n)
		}
	}
}

func (t *RBTree) removeBlackNode(n *node) {
	t.case1(n)
	t.removeLeaf(n)
}

func (t *RBTree) removeLeaf(leaf *node) {
	if leaf.key >= leaf.parent.key {
		leaf.parent.right = emptyLeaf
	} else {
		leaf.parent.left = emptyLeaf
	}
}

func (t *RBTree) leftRotation(n, parent, grandFather *node, recolor bool) {
	grandGrandFather := grandFather.parent
	t.updateParent(parent, grandFather, grandGrandFather)
	oldLeft := parent.left
	parent.left = grandFather
	grandFather.parent = parent

	grandFather.right = oldLeft
	oldLeft.parent = grandFather

	if recolor {
		parent.color = blackNode
		n.color = redNode
		grandFather.color = redNode
	}
}

func (t *RBTree) rightRotation(n, parent, grandFather *node, recolor bool) {
	grandGrandFather := grandFather.parent
	t.updateParent(parent, grandFather, grandGrandFather)
	oldRight := parent.right
	parent.right = grandFather
	grandFather.parent = parent

	grandFather.left = oldRight
	oldRight.parent = grandFather

	if recolor {
		parent.color = blackNode
		n.color = redNode
		grandFather.color = redNode
	}
}

func (t *RBTree) updateParent(n, parentOldChild, newParent *node) {
	n.parent = newParent
	if newParent != nil {
		if newParent.key > parentOldChild.key {
			newParent.left = n
		} else {
			newParent.right = n
		}
	} else {
		t.root = n
	}
}

func (t *RBTree) reColor(grandFather *node) {
	grandFather.right.color = blackNode
	grandFather.left.color = blackNode
	if !grandFather.equals(t.root) {
		grandFather.color = redNode
	}
	t.tryBalance(grandFather)
}

func (t *RBTree) case1(n *node) {
	/*
		Case 1 is when there's a double black node on the root
		Because we're at the root, we can simply remove it
		and reduce the black height of the whole tree.
		 __|10B|__                  __10B__
		/         \      ==>       /       \
		9B         20B            9B        20B
	*/
	if t.root.equals(n) {
		n.color = blackNode
		return
	}
	t.case2(n)
}

func (t *RBTree) case2(n *node) {
	/*
		Case 2 applies when
		the parent is BLACK
		the sibling is RED
		the sibling's children are BLACK or NIL
		It takes the sibling and rotates it
							 40B                                              60B
							/   \       --CASE 2 ROTATE-->                   /   \
						 |20B|   60R       LEFT ROTATE                      40R   80B
		DBL BLACK IS 20----^   /   \      SIBLING 60R                     /   \
							 50B    80B                                |20B|  50B
		(if the sibling's direction was left of it's parent, we would RIGHT ROTATE it)
		Now the original node's parent is RED, and we can apply case 4 or case 6
	*/
	parent := n.parent
	sibling, direction := t.getSibling(n)
	if sibling.color == redNode && parent.color == blackNode && sibling.left.color != redNode && sibling.right.color != redNode {
		if direction == leftDir {
			t.leftRotation(nil, sibling, parent, false)
		} else {
			t.rightRotation(nil, sibling, parent, false)
		}
		parent.color = redNode
		sibling.color = blackNode
		t.case1(n)
		return
	}
	t.case3(n)
}

func (t *RBTree) case3(n *node) {
	/*
		Case 3 deletion is when:
		the parent is BLACK
		the sibling is BLACK
		the sibling's children are BLACK
		Then, we make the sibling red and
		pass the double black node upwards
		                       Parent is black
		          ___50B___    Sibling is black                       ___50B___
		         /         \   Sibling's children are black          /         \
			   30B          80B        CASE 3                       30B        |80B|  Continue with other cases
		      /   \        /   \        ==>                        /  \        /   \
		    20B   35R    70B   |90B|<---REMOVE                   20B  35R     70R   X
		   /  \                                               /   \
		 34B   37B                                          34B   37B
	*/
	parent := n.parent
	sibling, _ := t.getSibling(n)
	if sibling.color == blackNode && parent.color == blackNode &&
		sibling.left.color != redNode && sibling.right.color != redNode {
		// color the sibling red and forward the double black node upwards
		// (call the cases again for the parent)
		sibling.color = redNode
		// start again
		t.case1(parent)
		return
	}
	t.case4(n)
}

func (t *RBTree) case4(n *node) {
	/*
		If the parent is red and the sibling is black with no red children,
		simply swap their colors
		DB-Double Black
			 __10R__                   __10B__         The black height of the left subtree has been incremented
		    /       \                 /       \        And the one below stays the same
		   DB        15B      ===>    X        15R     No consequences, we're done!
		            /   \                     /   \
		          12B   17B                 12B   17B
	*/
	parent := n.parent
	if parent.color == redNode {
		sibling, _ := t.getSibling(n)
		if sibling.color == blackNode && sibling.left.color != redNode &&
			sibling.right.color != redNode {
			// switch colors
			parent.color, sibling.color = sibling.color, parent.color
			return // Terminating
		}
	}
	t.case5(n)
}

func (t *RBTree) case5(n *node) {
	/*
			Case 5 is a rotation that changes the circumstances so that we can do a case 6
			If the closer node is red and the outer BLACK or NIL, we do a left/right rotation, depending on the orientation
			This will showcase when the CLOSER NODE's direction is RIGHT
				       ___50B___                                                    __50B__
		    	      /         \                                                  /       \
			        30B        |80B|  <-- Double black                           35B      |80B|        Case 6 is now
		    	   /  \        /   \      Closer node is red (35R)              /   \      /           applicable here,
			     20B  35R     70R   X     Outer is black (20B)               30R    37B  70R           so we redirect the node
		        	/   \                So we do a LEFT ROTATION          /   \                       to it :)
			       34B  37B               on 35R (closer node)           20B   34B
	*/
	sibling, direction := t.getSibling(n)
	closerNode := sibling.left
	if direction == leftDir {
		closerNode = sibling.right
	}
	outerNode := sibling.right
	if direction == leftDir {
		outerNode = sibling.left
	}
	if closerNode.color == redNode && outerNode.color != redNode && sibling.color == blackNode {
		if direction == leftDir {
			t.leftRotation(nil, closerNode, sibling, false)
		} else {
			t.rightRotation(nil, closerNode, sibling, false)
		}
		closerNode.color = blackNode
		sibling.color = redNode
	}
	t.case6(n)
}

func (t *RBTree) case6(n *node) {
	/*
		Case 6 requires
		    SIBLING to be BLACK
		    OUTER NODE to be RED
		Then, does a right/left rotation on the sibling
		This will showcase when the SIBLING's direction is LEFT
							Double Black
							__50B__       |                               __35B__
		                   /       \      |                              /       \
		     SIBLING--> 35B      |80B| <-                             30R       50R
		               /   \      /                                  /   \     /   \
		             30R    37B  70R   Outer node is RED            20B   34B 37B    80B
		            /   \              Closer node doesn't                           /
		          20B   34B                 matter                                   70R
		                                Parent doesn't
		                                    matter
		                          So we do a right rotation on 35B!
	*/
	sibling, direction := t.getSibling(n)
	outerNode := sibling.right
	if direction == leftDir {
		outerNode = sibling.left
	}

	if sibling.color == blackNode && outerNode.color == redNode {
		// terminating
		t.case6Rotation(sibling, direction)
		return
	}

	panic("We should have ended here, something is wrong")
}

func (t *RBTree) case6Rotation(sibling *node, direction dir) {
	parentColor := sibling.parent.color
	if direction == leftDir {
		t.leftRotation(nil, sibling, sibling.parent, false)
	} else {
		t.rightRotation(nil, sibling, sibling.parent, false)
	}
	// new parent is sibling
	sibling.color = parentColor
	sibling.right.color = blackNode
	sibling.left.color = blackNode
}

func (t *RBTree) tryBalance(n *node) {
	parent := n.parent
	key := n.key
	if parent == nil || // n is root
		parent.parent == nil || // parent is root
		(n.color != redNode || parent.color != redNode) {
		return
	}

	grandFather := parent.parent
	nodeDir := rightDir
	if parent.key > key {
		nodeDir = leftDir
	}
	parentDir := rightDir
	if grandFather.key > parent.key {
		parentDir = leftDir
	}
	uncle := grandFather.left
	if parentDir == leftDir {
		uncle = grandFather.right
	}
	generalDirection := nodeDir + parentDir
	if uncle.equals(emptyLeaf) || uncle.color == blackNode {
		// rotate
		if generalDirection == LLDir {
			t.rightRotation(n, parent, grandFather, true)
		} else if generalDirection == RRDir {
			t.leftRotation(n, parent, grandFather, true)
		} else if generalDirection == LRDir {
			t.rightRotation(nil, n, parent, false)
			// due to the prev rotation, our node is now the parent
			t.leftRotation(parent, n, grandFather, true)
		} else if generalDirection == RLDir {
			t.leftRotation(nil, n, parent, false)
			// # due to the prev rotation, our node is now the parent
			t.rightRotation(parent, n, grandFather, true)
		} else {
			panic(fmt.Sprintf("%s is not a valid direction!", generalDirection))
		}
	} else {
		// uncle is RED
		t.reColor(grandFather)
	}
}

func (t *RBTree) findParent(key string) (*node, dir) {
	return t.findNodeParent(t.root, key)
}

func (t *RBTree) findNodeParent(parent *node, key string) (*node, dir) {
	if key == parent.key {
		return parent, nilDir
	}
	if parent.key < key {
		if parent.right.color == nilNode {
			return parent, rightDir
		}
		return t.findNodeParent(parent.right, key)
	}
	if parent.left.color == nilNode {
		return parent, leftDir
	}
	return t.findNodeParent(parent.left, key)
}

func (t *RBTree) findNode(key string) *node {
	return t.findInnerNode(t.root, key)
}

func (t *RBTree) findInnerNode(n *node, key string) *node {
	if n == nil || n.color == nilNode {
		return nil
	}
	if key > n.key {
		return t.findInnerNode(n.right, key)
	}
	if key < n.key {
		return t.findInnerNode(n.left, key)
	}
	return n
}

func (t *RBTree) findInOrderSuccessor(n *node) *node {
	rNode := n.right
	lNode := rNode.left
	if lNode.equals(emptyLeaf) {
		return rNode
	}
	for !lNode.left.equals(emptyLeaf) {
		lNode = lNode.left
	}
	return lNode
}

func (t *RBTree) getSibling(n *node) (*node, dir) {
	/*
		e.g
			 20 (A)
			/     \
		 15(B)    25(C)

		getSibling(25(C)) => 15(B), 'R'
	*/
	parent := n.parent
	if n.key >= parent.key { // 证明我在右边，所以需要返回左边
		return parent.left, leftDir
	}
	return parent.right, rightDir
}

type nodeColor int

const (
	blackNode nodeColor = iota + 1
	redNode
	nilNode
)

func (n nodeColor) String() string {
	switch n {
	case blackNode:
		return "black"
	case redNode:
		return "red"
	case nilNode:
		return "nil"
	}
	return ""
}

type node struct {
	key    string
	value  any
	color  nodeColor
	parent *node
	left   *node
	right  *node
}

var (
	emptyLeaf = newNode("", nil, nilNode, nil)
)

func (n *node) String() string {
	return fmt.Sprintf("%s %v %s Node", n.key, n.value, n.color)
}

func (n *node) equals(other *node) bool {
	if n.color == nilNode && n.color == other.color {
		return true
	}

	parentsAreSame := false
	if n.parent == nil || other.parent == nil {
		parentsAreSame = n.parent == nil && other.parent == nil
	} else {
		parentsAreSame = n.parent.key == other.parent.key &&
			n.parent.color == other.parent.color
	}

	return n.key == other.key && n.color == other.color && parentsAreSame
}

func (n *node) hasChildren() bool {
	return n.getChildrenCount() > 0
}

func (n *node) getChildrenCount() int {
	if n.color == nilNode {
		return 0
	}

	sum := 0
	if n.left.color != nilNode {
		sum += 1
	}
	if n.right.color != nilNode {
		sum += 1
	}

	return sum
}

type nodeOption func(*node)

func right(right *node) nodeOption {
	return func(n *node) {
		n.right = right
	}
}

func left(left *node) nodeOption {
	return func(n *node) {
		n.left = left
	}
}

func newNode(key string, value any, color nodeColor,
	parent *node, options ...nodeOption) *node {
	n := &node{
		key:    key,
		value:  value,
		color:  color,
		parent: parent,
	}

	for _, option := range options {
		option(n)
	}

	return n
}
